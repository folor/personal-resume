<template>
  <!-- 3D 场景容器，position:relative 用于叠加 CSS2D 文字层 -->
  <div
    ref="containerRef"
    style="width: 100%; height: 650px; border: 1px solid #333; position: relative"
  ></div>
</template>

<script setup lang="ts">
// @ts-nocheck
import { ref, onMounted, onUnmounted } from "vue";
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
import { RGBELoader } from "three/examples/jsm/loaders/RGBELoader.js";
// 后期处理
import { EffectComposer } from "three/examples/jsm/postprocessing/EffectComposer.js";
import { RenderPass } from "three/examples/jsm/postprocessing/RenderPass.js";
import { OutlinePass } from "three/examples/jsm/postprocessing/OutlinePass.js";
import { OutputPass } from "three/examples/jsm/postprocessing/OutputPass.js";
import { UnrealBloomPass } from "three/examples/jsm/postprocessing/UnrealBloomPass.js";
import { SSAOPass } from "three/examples/jsm/postprocessing/SSAOPass.js";
// 文字标签渲染器
import {
  CSS2DRenderer,
  CSS2DObject,
} from "three/examples/jsm/renderers/CSS2DRenderer.js";

// ========== 基础场景引用 ==========
const containerRef = ref<HTMLDivElement | null>(null);

let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let labelRenderer: CSS2DRenderer;
let controls: OrbitControls;
let animationId: number | null = null;
const clock = new THREE.Clock();

let composer: EffectComposer;
let outlinePass: OutlinePass;
let bloomPass: UnrealBloomPass;
let ssaoPass: SSAOPass;

// ========== 模型相关 ==========
interface ModelWrap {
  group: THREE.Group;
  mixer: THREE.AnimationMixer;
  runAction?: THREE.AnimationAction;
  idleAction?: THREE.AnimationAction;
}
const modelWrapList: ModelWrap[] = [];
const modelList: THREE.Group[] = [];

let selectedModel: THREE.Group | null = null;
const originMaterialMap = new WeakMap<THREE.Mesh, THREE.Material | THREE.Material[]>();

// ========== 头顶文字标签 ==========
const labelMap = new Map<THREE.Group, CSS2DObject>();
const labelTimeoutMap = new Map<THREE.Group, number>();
const LABEL_DURATION = 1500;

// ========== 弹跳动画队列 ==========
interface ModelAnimItem {
  model: THREE.Group;
  originY: number;
  height: number;
  originZ: number;
  moveDistance: number;
  progress: number;
}
let modelAnimQueue: ModelAnimItem[] = [];

// ========== 自动移动 ==========
interface MoveTarget {
  x: number;
  z: number;
  pickupItem?: THREE.Mesh | null;
}
const moveTargetMap = new Map<THREE.Group, MoveTarget>();
const MOVE_SPEED = 3.0;

// ========== 呼吸缩放 ==========
const BASE_SCALE = 3.0;
const BREATH_AMPLITUDE = 0.05;
const BREATH_SPEED = 2.0;

// ========== 鼠标交互 ==========
let mouseDownTime = 0;
let isDragMode = false;
let downHitModel: THREE.Group | null = null;
let downHitPoint: THREE.Vector3 | null = null;
let downGroundPoint: THREE.Vector3 | null = null;
let downHitItem: THREE.Mesh | null = null;
const LONG_PRESS_MS = 300;
const dragOffset = new THREE.Vector2();

// ========== 可拾取物品 ==========
let pickupItemList: THREE.Mesh[] = [];

// ========== 点击光圈特效 ==========
interface ClickRipple {
  mesh: THREE.Mesh;
  targetModel: THREE.Group;
  fading: boolean;
  fadeOpacity: number;
}
let rippleList: ClickRipple[] = [];

// ========== 爆炸粒子+拖尾 ==========
interface ExplodeParticleItem {
  point: THREE.Points;
  trailLine: THREE.Line;
  velocity: THREE.Vector3;
  life: number;
  maxLife: number;
  active: boolean;
  trailLength: number;
  trailPositions: Float32Array;
  pointMat: THREE.PointsMaterial;
  trailMat: THREE.LineBasicMaterial;
}
let particlePool: ExplodeParticleItem[] = [];
const POOL_SIZE = 220;
let prevPointerPos = new THREE.Vector2();
let pointerSpeed = 0;

function createOneParticle(): ExplodeParticleItem {
  const geo = new THREE.BufferGeometry();
  const posArr = new Float32Array([0, 0, 0]);
  geo.setAttribute("position", new THREE.BufferAttribute(posArr, 3));
  const pMat = new THREE.PointsMaterial({
    size: 0.22,
    transparent: true,
    opacity: 0.9,
    sizeAttenuation: true,
    depthWrite: false,
  });
  const point = new THREE.Points(geo, pMat);

  const trailLen = 12;
  const trailArr = new Float32Array(trailLen * 3);
  const trailGeo = new THREE.BufferGeometry();
  trailGeo.setAttribute("position", new THREE.BufferAttribute(trailArr, 3));
  const tMat = new THREE.LineBasicMaterial({
    transparent: true,
    opacity: 0.45,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
  });
  const trailLine = new THREE.Line(trailGeo, tMat);

  return {
    point,
    trailLine,
    velocity: new THREE.Vector3(),
    life: 0,
    maxLife: 110,
    active: false,
    trailLength: trailLen,
    trailPositions: trailArr,
    pointMat: pMat,
    trailMat: tMat,
  };
}

function initParticlePool() {
  for (let i = 0; i < POOL_SIZE; i++) {
    const p = createOneParticle();
    scene.add(p.point);
    scene.add(p.trailLine);
    particlePool.push(p);
  }
}

function spawnExplode(pos: THREE.Vector3, speedVal: number) {
  const burstCount = 45;
  const t = THREE.MathUtils.clamp(speedVal, 0, 1);
  const burstColor = new THREE.Color().setHSL(0.55 - t * 0.25, 0.78, 0.62);

  let spawned = 0;
  for (const p of particlePool) {
    if (spawned >= burstCount) break;
    if (!p.active) {
      p.active = true;
      p.life = p.maxLife;
      p.point.position.copy(pos);
      const dir = new THREE.Vector3(
        (Math.random() - 0.5) * 2,
        (Math.random() - 0.5) * 2,
        (Math.random() - 0.5) * 2
      ).normalize();
      const spd = 0.25 + Math.random() * 0.65;
      p.velocity.copy(dir).multiplyScalar(spd);

      p.pointMat.color.copy(burstColor);
      p.trailMat.color.copy(burstColor);

      for (let k = 0; k < p.trailLength; k++) {
        const idx = k * 3;
        p.trailPositions[idx] = pos.x;
        p.trailPositions[idx + 1] = pos.y;
        p.trailPositions[idx + 2] = pos.z;
      }
      p.trailLine.geometry.attributes.position.needsUpdate = true;
      spawned++;
    }
  }
}

function tickExplodeParticles() {
  for (const p of particlePool) {
    if (!p.active) continue;
    p.life -= 1;
    if (p.life <= 0) {
      p.active = false;
      continue;
    }
    p.point.position.add(p.velocity);
    p.velocity.multiplyScalar(0.97);

    for (let i = p.trailLength - 1; i > 0; i--) {
      const cur = i * 3;
      const prev = (i - 1) * 3;
      p.trailPositions[cur] = p.trailPositions[prev];
      p.trailPositions[cur + 1] = p.trailPositions[prev + 1];
      p.trailPositions[cur + 2] = p.trailPositions[prev + 2];
    }
    p.trailPositions[0] = p.point.position.x;
    p.trailPositions[1] = p.point.position.y;
    p.trailPositions[2] = p.point.position.z;
    p.trailLine.geometry.attributes.position.needsUpdate = true;

    const alpha = p.life / p.maxLife;
    p.pointMat.opacity = alpha;
    p.trailMat.opacity = alpha * 0.5;
  }
}

function disposeParticlePool() {
  for (const p of particlePool) {
    scene.remove(p.point);
    scene.remove(p.trailLine);
    p.point.geometry.dispose();
    p.pointMat.dispose();
    p.trailLine.geometry.dispose();
    p.trailMat.dispose();
  }
  particlePool = [];
}

function calcPointerSpeed(clientX: number, clientY: number) {
  const dx = clientX - prevPointerPos.x;
  const dy = clientY - prevPointerPos.y;
  pointerSpeed = Math.sqrt(dx * dx + dy * dy) / 130;
  prevPointerPos.set(clientX, clientY);
}

// ========== 工具函数 ==========
function removeRippleByModel(model: THREE.Group) {
  for (let i = rippleList.length - 1; i >= 0; i--) {
    if (rippleList[i].targetModel === model) {
      const item = rippleList[i];
      scene.remove(item.mesh);
      item.mesh.geometry.dispose();
      item.mesh.material.dispose();
      rippleList.splice(i, 1);
    }
  }
}

function createClickRipple(x: number, z: number, targetModel: THREE.Group) {
  removeRippleByModel(targetModel);
  const geo = new THREE.RingGeometry(0.3, 0.5, 48);
  const mat = new THREE.MeshBasicMaterial({
    color: 0x4da6ff,
    transparent: true,
    opacity: 0.65,
    side: THREE.DoubleSide,
    depthWrite: false,
  });
  const ring = new THREE.Mesh(geo, mat);
  ring.rotation.x = -Math.PI / 2;
  ring.position.set(x, 0.02, z);
  scene.add(ring);
  rippleList.push({ mesh: ring, targetModel, fading: false, fadeOpacity: 0.65 });
}

function createPickupItems() {
  const itemConfigs = [
    { pos: new THREE.Vector3(5, 0.5, 3), color: 0xffd700, emissive: 0xffaa00 },
    { pos: new THREE.Vector3(-4, 0.5, 4), color: 0x00ff88, emissive: 0x00cc66 },
    { pos: new THREE.Vector3(2, 0.5, -5), color: 0xff66cc, emissive: 0xcc3399 },
    { pos: new THREE.Vector3(6, 0.5, -2), color: 0x66aaff, emissive: 0x3377dd },
  ];

  itemConfigs.forEach((cfg) => {
    const geo = new THREE.SphereGeometry(0.35, 32, 32);
    const mat = new THREE.MeshStandardMaterial({
      color: cfg.color,
      emissive: cfg.emissive,
      emissiveIntensity: 0.7,
      metalness: 0.9,
      roughness: 0.15,
    });
    const item = new THREE.Mesh(geo, mat);
    item.position.copy(cfg.pos);
    item.castShadow = true;
    scene.add(item);
    pickupItemList.push(item);
  });
}

function removePickupItem(item: THREE.Mesh) {
  const idx = pickupItemList.indexOf(item);
  if (idx > -1) pickupItemList.splice(idx, 1);
  scene.remove(item);
  item.geometry.dispose();
  item.material.dispose();
}

function createModelLabel(
  model: THREE.Group,
  text: string,
  duration: number = LABEL_DURATION
) {
  removeModelLabel(model);

  const div = document.createElement("div");
  div.textContent = text;
  div.style.cssText = `
    padding: 4px 12px;
    background: rgba(245, 250, 255, 0.88);
    color: #173b5c;
    border-radius: 6px;
    font-size: 14px;
    font-family: sans-serif;
    white-space: nowrap;
    pointer-events: none;
    border: 1px solid rgba(77, 166, 255, 0.7);
    box-shadow: 0 0 12px rgba(77, 166, 255, 0.35);
    transform: translateY(-4px);
  `;
  const label = new CSS2DObject(div);
  label.position.set(0, 2.2, 0);
  model.add(label);
  labelMap.set(model, label);

  const timer = window.setTimeout(() => {
    removeModelLabel(model);
  }, duration);
  labelTimeoutMap.set(model, timer);
}

function removeModelLabel(model: THREE.Group) {
  const timer = labelTimeoutMap.get(model);
  if (timer) {
    clearTimeout(timer);
    labelTimeoutMap.delete(model);
  }
  const label = labelMap.get(model);
  if (label) {
    model.remove(label);
    labelMap.delete(model);
  }
}

function getModelWrap(group: THREE.Group): ModelWrap | undefined {
  return modelWrapList.find((item) => item.group === group);
}

function playRun(model: THREE.Group, play: boolean) {
  const wrap = getModelWrap(model);
  if (!wrap?.runAction) return;
  if (play) {
    wrap.runAction.play();
    wrap.runAction.setLoop(THREE.LoopRepeat, Infinity);
  } else {
    wrap.runAction.stop();
  }
}

// ========== 射线拾取 ==========
const raycaster = new THREE.Raycaster();
const mouse = new THREE.Vector2();
const dragPlane = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0);
let isDragging = false;
const tempWorldPoint = new THREE.Vector3();

function onMouseDown(e: MouseEvent) {
  if (!renderer.domElement || modelList.length === 0) return;
  const rect = renderer.domElement.getBoundingClientRect();
  mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
  mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;

  calcPointerSpeed(e.clientX, e.clientY);

  raycaster.setFromCamera(mouse, camera);
  if (raycaster.ray.intersectPlane(dragPlane, tempWorldPoint)) {
    spawnExplode(tempWorldPoint, pointerSpeed);
  }

  mouseDownTime = performance.now();
  isDragMode = false;
  isDragging = false;
  downHitItem = null;

  raycaster.setFromCamera(mouse, camera);

  for (const item of pickupItemList) {
    const hits = raycaster.intersectObject(item, false);
    if (hits.length > 0) {
      downHitItem = item;
      break;
    }
  }

  const hitList = raycaster.intersectObjects(modelList, true);
  if (hitList.length > 0) {
    let targetObj = hitList[0].object;
    while (targetObj.parent && !modelList.includes(targetObj as THREE.Group)) {
      targetObj = targetObj.parent;
    }
    downHitModel = targetObj as THREE.Group;
    downHitPoint = hitList[0].point.clone();
  } else {
    downHitModel = null;
    downHitPoint = null;
  }

  const groundPoint = new THREE.Vector3();
  if (raycaster.ray.intersectPlane(dragPlane, groundPoint)) {
    downGroundPoint = groundPoint.clone();
  } else {
    downGroundPoint = null;
  }
}

function onMouseMove(e: MouseEvent) {
  if (!renderer.domElement) return;
  const rect = renderer.domElement.getBoundingClientRect();
  mouse.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
  mouse.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
  calcPointerSpeed(e.clientX, e.clientY);
  raycaster.setFromCamera(mouse, camera);

  if (!isDragMode) {
    const elapsed = performance.now() - mouseDownTime;
    if (elapsed > LONG_PRESS_MS && downHitModel && !downHitItem) {
      isDragMode = true;
      isDragging = true;
      selectedModel = downHitModel;
      outlinePass.selectedObjects = [selectedModel];
      moveTargetMap.delete(selectedModel);
      removeRippleByModel(selectedModel);
      playRun(selectedModel, false);
      controls.enabled = false;
      dragOffset.x = selectedModel.position.x - downHitPoint!.x;
      dragOffset.y = selectedModel.position.z - downHitPoint!.z;
    }
  }

  if (isDragging && selectedModel) {
    raycaster.setFromCamera(mouse, camera);
    const intersectPoint = new THREE.Vector3();
    raycaster.ray.intersectPlane(dragPlane, intersectPoint);
    selectedModel.position.x = intersectPoint.x + dragOffset.x;
    selectedModel.position.z = intersectPoint.z + dragOffset.y;
    if (selectedModel.position.y < 0) selectedModel.position.y = 0;
  }

  const SPEED_THRESHOLD = 0.02;
  if (pointerSpeed > SPEED_THRESHOLD && !isDragMode) {
    const worldPos = new THREE.Vector3();
    if (raycaster.ray.intersectPlane(dragPlane, worldPos)) {
      spawnExplode(worldPos, pointerSpeed);
    }
  }
}

function onMouseUp(e: MouseEvent) {
  const elapsed = performance.now() - mouseDownTime;

  if (isDragMode) {
    isDragging = false;
    isDragMode = false;
    controls.enabled = true;
    downHitModel = null;
    downHitPoint = null;
    downGroundPoint = null;
    downHitItem = null;
    return;
  }

  if (elapsed < LONG_PRESS_MS) {
    if (downHitItem) {
      const targetModel = selectedModel || modelList[0];
      if (targetModel) {
        if (!selectedModel) {
          selectedModel = targetModel;
          outlinePass.selectedObjects = [selectedModel];
        }
        createModelLabel(targetModel, "去捡东西");
        createClickRipple(downHitItem.position.x, downHitItem.position.z, targetModel);
        moveTargetMap.set(targetModel, {
          x: downHitItem.position.x,
          z: downHitItem.position.z,
          pickupItem: downHitItem,
        });
        playRun(targetModel, true);
      }
    } else if (downHitModel) {
      if (selectedModel === downHitModel) {
        createModelLabel(selectedModel, "你好~");
        modelAnimQueue.push({
          model: selectedModel,
          originY: Math.max(0, selectedModel.position.y),
          height: 1.8,
          originZ: selectedModel.position.z,
          moveDistance: 1.0,
          progress: 0,
        });
      } else {
        clearHighlight();
        if (selectedModel) removeModelLabel(selectedModel);
        selectedModel = downHitModel;
        outlinePass.selectedObjects = [selectedModel];
        moveTargetMap.delete(selectedModel);
        removeRippleByModel(selectedModel);
        playRun(selectedModel, false);
        createModelLabel(selectedModel, "你好~");
        modelAnimQueue.push({
          model: selectedModel,
          originY: Math.max(0, selectedModel.position.y),
          height: 1.8,
          originZ: selectedModel.position.z,
          moveDistance: 1.0,
          progress: 0,
        });
      }
    } else if (downGroundPoint) {
      const targetModel = selectedModel || modelList[0];
      if (targetModel) {
        if (!selectedModel) {
          selectedModel = targetModel;
          outlinePass.selectedObjects = [selectedModel];
        }
        createModelLabel(targetModel, "我来啦！！！");
        createClickRipple(downGroundPoint.x, downGroundPoint.z, targetModel);
        moveTargetMap.set(targetModel, {
          x: downGroundPoint.x,
          z: downGroundPoint.z,
        });
        playRun(targetModel, true);
      }
    }
  }

  downHitModel = null;
  downHitPoint = null;
  downGroundPoint = null;
  downHitItem = null;
  isDragMode = false;
  isDragging = false;
  controls.enabled = true;
}

// ========== 触屏事件兼容 ==========
function onTouchStart(e: TouchEvent) {
  e.preventDefault();
  const t = e.touches[0];
  const mouseEvt = new MouseEvent("mousedown", {
    clientX: t.clientX,
    clientY: t.clientY,
  });
  onMouseDown(mouseEvt);
}
function onTouchMove(e: TouchEvent) {
  e.preventDefault();
  const t = e.touches[0];
  const mouseEvt = new MouseEvent("mousemove", {
    clientX: t.clientX,
    clientY: t.clientY,
  });
  onMouseMove(mouseEvt);
}
function onTouchEnd(e: TouchEvent) {
  e.preventDefault();
  const mouseEvt = new MouseEvent("mouseup", { clientX: 0, clientY: 0 });
  onMouseUp(mouseEvt);
}

function setHighlight(obj: THREE.Group) {
  obj.traverse((mesh) => {
    const m = mesh as THREE.Mesh;
    if (!m.isMesh) return;
    if (!originMaterialMap.has(m)) {
      originMaterialMap.set(m, m.material);
    }
    const originMat = m.material as THREE.MeshStandardMaterial;
    const highlightMat = new THREE.MeshStandardMaterial({
      color: originMat.color,
      map: originMat.map,
      emissive: 0xffdd00,
      emissiveIntensity: 0.35,
    });
    m.material = highlightMat;
  });
}

function clearHighlight() {
  modelList.forEach((group) => {
    group.traverse((mesh) => {
      const m = mesh as THREE.Mesh;
      if (!m.isMesh) return;
      if (originMaterialMap.has(m)) {
        m.material = originMaterialMap.get(m)!;
      }
    });
  });
}

// ========== 加载HDRI环境贴图 ==========
function loadHDRI() {
  const pmremGenerator = new THREE.PMREMGenerator(renderer);
  pmremGenerator.compileEquirectangularShader();

  const rgbeLoader = new RGBELoader();
  // 科技感室内HDR，免费可商用
  rgbeLoader.load(
    "https://threejs.org/examples/textures/equirectangular/royal_esplanade_1k.hdr",
    (texture) => {
      const envMap = pmremGenerator.fromEquirectangular(texture).texture;
      scene.environment = envMap;
      texture.dispose();
      pmremGenerator.dispose();
      console.log("HDRI环境贴图加载完成");
    },
    undefined,
    (err) => {
      console.warn("HDRI加载失败，使用默认环境:", err);
    }
  );
}

// ========== 初始化场景 ==========
function init() {
  if (!containerRef.value) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;

  // 场景：蓝白渐变 + 柔和雾化
  scene = new THREE.Scene();
  scene.background = new THREE.Color(0xe6f4ff);
  scene.fog = new THREE.FogExp2(0xeaf4ff, 0.012);

  // 相机
  camera = new THREE.PerspectiveCamera(75, w / h, 0.1, 2000);
  camera.position.set(12, 12, 12);

  // WebGL 渲染器：开启色调映射，配合Bloom效果更好
  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(w, h);
  renderer.setPixelRatio(window.devicePixelRatio);
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  renderer.toneMapping = THREE.ACESFilmicToneMapping;
  renderer.toneMappingExposure = 1.0;
  containerRef.value.appendChild(renderer.domElement);

  // 文字标签渲染器
  labelRenderer = new CSS2DRenderer();
  labelRenderer.setSize(w, h);
  labelRenderer.domElement.style.position = "absolute";
  labelRenderer.domElement.style.top = "0";
  labelRenderer.domElement.style.left = "0";
  labelRenderer.domElement.style.pointerEvents = "none";
  containerRef.value.appendChild(labelRenderer.domElement);

  // 灯光
  const dirLight = new THREE.DirectionalLight(0xffffff, 1.3);
  dirLight.position.set(12, 22, 12);
  dirLight.castShadow = true;
  dirLight.shadow.mapSize.set(2048, 2048);
  dirLight.shadow.camera.near = 0.5;
  dirLight.shadow.camera.far = 100;
  dirLight.shadow.camera.left = -30;
  dirLight.shadow.camera.right = 30;
  dirLight.shadow.camera.top = 30;
  dirLight.shadow.camera.bottom = -30;
  dirLight.shadow.bias = -0.0001;
  scene.add(dirLight);

  // 半球光：模拟天空+地面反射，替代HDRI，不需要外部资源
  const hemiLight = new THREE.HemisphereLight(0xb8dcff, 0xd7ecff, 0.7);
  scene.add(hemiLight);

  const rimLight = new THREE.DirectionalLight(0x4da6ff, 0.4);
  rimLight.position.set(-10, 8, 10);
  scene.add(rimLight);

  // 地面
  const groundGeo = new THREE.PlaneGeometry(200, 200);
  const groundMat = new THREE.MeshStandardMaterial({
    color: 0xd7ecff,
    roughness: 0.45,
    metalness: 0.35,
  });
  const ground = new THREE.Mesh(groundGeo, groundMat);
  ground.rotation.x = -Math.PI / 2;
  ground.receiveShadow = true;
  scene.add(ground);

  // 网格
  const grid = new THREE.GridHelper(25, 25, 0x4da6ff, 0x91c8ff);
  grid.position.y = 0.02;
  scene.add(grid);

  const gridFine = new THREE.GridHelper(50, 50, 0xb8dcff, 0xb8dcff);
  gridFine.position.y = 0.015;
  scene.add(gridFine);

  const axes = new THREE.AxesHelper(6);
  axes.material.opacity = 0.75;
  axes.material.transparent = true;
  scene.add(axes);

  // 轨道控制器
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.05;
  controls.minDistance = 2;
  controls.maxDistance = 120;

  // 鼠标事件
  renderer.domElement.addEventListener("mousedown", onMouseDown);
  renderer.domElement.addEventListener("mousemove", onMouseMove);
  renderer.domElement.addEventListener("mouseup", onMouseUp);
  renderer.domElement.addEventListener("mouseleave", onMouseUp);

  // 触屏事件
  renderer.domElement.addEventListener("touchstart", onTouchStart, { passive: false });
  renderer.domElement.addEventListener("touchmove", onTouchMove, { passive: false });
  renderer.domElement.addEventListener("touchend", onTouchEnd, { passive: false });

  // ========== 后期处理管线 ==========
  // 顺序：RenderPass → SSAO → Outline → Bloom → Output
  composer = new EffectComposer(renderer);

  // 1. 基础渲染
  composer.addPass(new RenderPass(scene, camera));

  // 2. SSAO 环境光遮蔽（物体接触阴影，增强立体感）
  ssaoPass = new SSAOPass(scene, camera, w, h);
  ssaoPass.kernelRadius = 8;
  ssaoPass.minDistance = 0.005;
  ssaoPass.maxDistance = 0.05;
  ssaoPass.output = SSAOPass.OUTPUT.Default;
  composer.addPass(ssaoPass);

  // 3. 选中边缘光晕
  outlinePass = new OutlinePass(new THREE.Vector2(w, h), scene, camera);
  outlinePass.edgeStrength = 3.0;
  outlinePass.edgeGlow = 0.65;
  outlinePass.edgeThickness = 1.2;
  outlinePass.pulsePeriod = 0;
  outlinePass.visibleEdgeColor.set(0x4da6ff);
  outlinePass.hiddenEdgeColor.set(0x000000);
  composer.addPass(outlinePass);

  // 4. Bloom 泛光（发光元素更炫）
  bloomPass = new UnrealBloomPass(
    new THREE.Vector2(w, h),
    0.35, // strength 泛光强度
    0.5, // radius 泛光半径
    0.85 // threshold 亮度阈值（超过这个亮度才发光）
  );
  composer.addPass(bloomPass);

  // 5. 输出
  composer.addPass(new OutputPass());

  // 加载HDRI环境贴图
  // loadHDRI();

  // 创建可拾取物品
  createPickupItems();

  // 初始化爆炸粒子池
  initParticlePool();

  // 加载 GLB 模型（带 run 动画）
  loadModelWithAnim("/scene.glb", new THREE.Vector3(-3, 0, 0));

  // ========== 动画循环 ==========
  function animate() {
    animationId = requestAnimationFrame(animate);
    controls.update();
    const delta = clock.getDelta();
    const elapsed = clock.getElapsedTime();

    // 1. 骨骼动画
    modelWrapList.forEach((wrap) => wrap.mixer.update(delta));

    // 2. 模型呼吸缩放
    const breathScale = BASE_SCALE + Math.sin(elapsed * BREATH_SPEED) * BREATH_AMPLITUDE;
    modelList.forEach((model) => model.scale.set(breathScale, breathScale, breathScale));

    // 3. 物品自转 + 浮动
    pickupItemList.forEach((item, i) => {
      item.rotation.y += delta * 2;
      item.position.y = 0.5 + Math.sin(elapsed * 3 + i * 1.2) * 0.15;
    });

    // 4. 爆炸粒子
    tickExplodeParticles();

    // 5. 自动移动
    moveTargetMap.forEach((target, model) => {
      const dx = target.x - model.position.x;
      const dz = target.z - model.position.z;
      const distance = Math.sqrt(dx * dx + dz * dz);
      if (distance < 0.05) {
        model.position.x = target.x;
        model.position.z = target.z;
        moveTargetMap.delete(model);
        playRun(model, false);

        if (target.pickupItem) {
          removePickupItem(target.pickupItem);
          createModelLabel(model, "捡到物品");
        }
      } else {
        const moveDistance = MOVE_SPEED * delta;
        const ratio = Math.min(moveDistance / distance, 1);
        model.position.x += dx * ratio;
        model.position.z += dz * ratio;
        model.rotation.y = Math.atan2(dx, dz);
      }
      if (model.position.y < 0) model.position.y = 0;
    });

    // 6. 弹跳动画
    const speed = 0.042;
    for (let i = modelAnimQueue.length - 1; i >= 0; i--) {
      const item = modelAnimQueue[i];
      item.progress += speed;
      if (item.progress >= 1) {
        item.progress = 1;
        modelAnimQueue.splice(i, 1);
      }
      const bounceOffset = Math.sin(Math.PI * item.progress) * item.height;
      item.model.position.y = item.originY + bounceOffset;
      item.model.position.z = THREE.MathUtils.lerp(
        item.originZ,
        item.originZ + item.moveDistance,
        item.progress
      );
    }

    // 7. 光圈特效
    for (let i = rippleList.length - 1; i >= 0; i--) {
      const item = rippleList[i];
      const stillMoving = moveTargetMap.has(item.targetModel);
      if (stillMoving && !item.fading) {
        const pulse = 1.0 + Math.sin(elapsed * 4) * 0.15;
        item.mesh.scale.set(pulse, pulse, pulse);
        item.mesh.material.opacity = 0.55 + Math.sin(elapsed * 4) * 0.1;
      } else {
        if (!item.fading) {
          item.fading = true;
          item.fadeOpacity = item.mesh.material.opacity;
        }
        item.fadeOpacity -= delta * 2.5;
        item.mesh.material.opacity = Math.max(0, item.fadeOpacity);
        item.mesh.scale.multiplyScalar(1 + delta * 0.8);
        if (item.fadeOpacity <= 0) {
          scene.remove(item.mesh);
          item.mesh.geometry.dispose();
          item.mesh.material.dispose();
          rippleList.splice(i, 1);
        }
      }
    }

    // 8. Y 轴钳制
    modelList.forEach((model) => {
      if (model.position.y < 0) model.position.y = 0;
    });

    // 渲染
    composer.render();
    labelRenderer.render(scene, camera);
  }
  animate();

  window.addEventListener("resize", onWindowResize);
}

function loadModelWithAnim(url: string, pos: THREE.Vector3) {
  const loader = new GLTFLoader();
  loader.load(
    url,
    (gltf) => {
      const model = gltf.scene;
      model.position.copy(pos);
      if (model.position.y < 0) model.position.y = 0;
      model.scale.set(BASE_SCALE, BASE_SCALE, BASE_SCALE);
      model.traverse((obj) => {
        obj.castShadow = true;
        obj.receiveShadow = true;
      });
      scene.add(model);
      modelList.push(model);

      const mixer = new THREE.AnimationMixer(model);
      let runAction: THREE.AnimationAction | undefined;
      const runClip = THREE.AnimationClip.findByName(gltf.animations, "run");
      if (runClip) {
        runAction = mixer.clipAction(runClip);
        runAction.setLoop(THREE.LoopRepeat, Infinity);
      } else {
        console.warn("⚠️未找到名为run的骨骼动画片段");
      }
      modelWrapList.push({ group: model, mixer, runAction });
    },
    (p) => console.log(`加载进度 ${Math.round((p.loaded / p.total) * 100)}%`),
    (err) => console.error("模型加载失败", err)
  );
}

function onWindowResize() {
  if (!containerRef.value) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;
  camera.aspect = w / h;
  camera.updateProjectionMatrix();
  renderer.setSize(w, h);
  composer.setSize(w, h);
  labelRenderer.setSize(w, h);
  // SSAO和Bloom也要同步尺寸
  if (ssaoPass) ssaoPass.setSize(w, h);
  if (bloomPass) bloomPass.setSize(w, h);
}

function dispose() {
  window.removeEventListener("resize", onWindowResize);
  if (animationId !== null) cancelAnimationFrame(animationId);

  renderer.domElement.removeEventListener("mousedown", onMouseDown);
  renderer.domElement.removeEventListener("mousemove", onMouseMove);
  renderer.domElement.removeEventListener("mouseup", onMouseUp);
  renderer.domElement.removeEventListener("mouseleave", onMouseUp);
  renderer.domElement.removeEventListener("touchstart", onTouchStart);
  renderer.domElement.removeEventListener("touchmove", onTouchMove);
  renderer.domElement.removeEventListener("touchend", onTouchEnd);

  controls?.dispose();
  composer?.dispose();
  renderer?.dispose();

  disposeParticlePool();

  labelTimeoutMap.forEach((timer) => clearTimeout(timer));
  labelTimeoutMap.clear();
  labelMap.forEach((label, model) => model.remove(label));
  labelMap.clear();
  if (labelRenderer?.domElement) {
    labelRenderer.domElement.remove();
  }

  pickupItemList.forEach((item) => {
    scene.remove(item);
    item.geometry.dispose();
    item.material.dispose();
  });
  pickupItemList = [];

  rippleList.forEach((item) => {
    scene.remove(item.mesh);
    item.mesh.geometry.dispose();
    item.mesh.material.dispose();
  });
  rippleList = [];

  modelAnimQueue = [];
  moveTargetMap.clear();
  modelWrapList.forEach((wrap) => wrap.mixer.stopAllAction());
  modelWrapList.length = 0;

  modelList.forEach((m) => {
    scene.remove(m);
    m.traverse((obj) => {
      if (obj.geometry) obj.geometry.dispose();
      if (obj.material) {
        if (Array.isArray(obj.material)) obj.material.forEach((s) => s.dispose());
        else obj.material.dispose();
      }
    });
  });
  modelList.length = 0;
}

onMounted(() => init());
onUnmounted(() => dispose());
</script>
