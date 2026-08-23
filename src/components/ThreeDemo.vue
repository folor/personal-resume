<template>
  <!-- 3D 场景容器，position:relative 用于叠加 CSS2D 文字层 / 工具栏 / HUD -->
  <div ref="containerRef" class="bot-demo">
    <div class="bot-toolbar">
      <button :class="{ active: followCam }" @click="followCam = !followCam">跟随相机</button>
      <button :class="{ active: wanderAI }" @click="wanderAI = !wanderAI">NPC 巡逻</button>
      <button :class="{ active: isNight }" @click="toggleNight">
        {{ isNight ? "☀️ 白天" : "🌙 夜晚" }}
      </button>
    </div>

    <div class="bot-hud">
      <span>⭐ 得分 <b>{{ score }}</b></span>
      <span>🔮 能量球 <b>{{ itemCount }}</b></span>
      <span>🕹️ WASD / 方向键移动</span>
    </div>

    <div v-if="loadingCount > 0" class="bot-loading">
      模型加载中 {{ 3 - loadingCount }}/3 …
    </div>
  </div>
</template>

<script setup lang="ts">
// @ts-nocheck
import { ref, onMounted, onUnmounted } from "vue";
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";
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
const score = ref(0);
const itemCount = ref(0);
const loadingCount = ref(3);
const followCam = ref(false);
const wanderAI = ref(true);
const isNight = ref(false);

let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let labelRenderer: CSS2DRenderer;
let controls: OrbitControls;
let animationId: number | null = null;
let resizeObserver: ResizeObserver | null = null;
let disposed = false;
const clock = new THREE.Clock();

let composer: EffectComposer;
let outlinePass: OutlinePass;
let bloomPass: UnrealBloomPass;
let ssaoPass: SSAOPass;
let dirLight: THREE.DirectionalLight;
let hemiLight: THREE.HemisphereLight;
let rimLight: THREE.DirectionalLight;

// ========== 模型相关 ==========
interface ModelWrap {
  group: THREE.Group;
  mixer: THREE.AnimationMixer;
  runAction?: THREE.AnimationAction;
  idleAction?: THREE.AnimationAction;
}
let modelWrapList: ModelWrap[] = [];
let modelList: THREE.Group[] = [];

let selectedModel: THREE.Group | null = null;

// ========== 头顶文字标签 ==========
const labelMap = new Map<THREE.Group, CSS2DObject>();
const labelTimeoutMap = new Map<THREE.Group, number>();
const nameTagMap = new Map<THREE.Group, CSS2DObject>();
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
const respawnTimers: number[] = [];
const ITEM_COLORS = [0xffd700, 0x00ff88, 0xff66cc, 0x66aaff, 0xff8844];
const ITEM_EMISSIVES = [0xffaa00, 0x00cc66, 0xcc3399, 0x3377dd, 0xcc5511];

// ========== 障碍物 ==========
interface Obstacle {
  x: number;
  z: number;
  r: number;
}
let obstacleList: Obstacle[] = [];
let obstacleMeshes: THREE.Mesh[] = [];
const BOUND = 12;

// ========== 键盘控制 ==========
const keys = new Set<string>();
let keyboardMoving = false;
const KEY_MAP: Record<string, [number, number]> = {
  w: [0, -1], arrowup: [0, -1],
  s: [0, 1], arrowdown: [0, 1],
  a: [-1, 0], arrowleft: [-1, 0],
  d: [1, 0], arrowright: [1, 0],
};

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

// ========== 昼夜主题 ==========
const DAY = {
  bg: 0xe6f4ff,
  fog: 0xeaf4ff,
  fogDensity: 0.012,
  dirColor: 0xffffff,
  dirIntensity: 1.3,
  hemiIntensity: 0.7,
  rimColor: 0x4da6ff,
  bloom: 0.35,
};
const NIGHT = {
  bg: 0x0a1226,
  fog: 0x0a1226,
  fogDensity: 0.02,
  dirColor: 0x8899ff,
  dirIntensity: 0.35,
  hemiIntensity: 0.22,
  rimColor: 0x3377ff,
  bloom: 0.75,
};

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

// ========== 能量球（拾取物品）==========
function randomItemPos(): THREE.Vector3 {
  // 避开障碍物附近
  for (let tryN = 0; tryN < 10; tryN++) {
    const x = (Math.random() * 2 - 1) * (BOUND - 1);
    const z = (Math.random() * 2 - 1) * (BOUND - 1);
    const tooClose = obstacleList.some((o) => Math.hypot(o.x - x, o.z - z) < o.r + 1.2);
    if (!tooClose) return new THREE.Vector3(x, 0.5, z);
  }
  return new THREE.Vector3(5, 0.5, 3);
}

function spawnOneItem(idx?: number) {
  const i = idx !== undefined ? idx % ITEM_COLORS.length : Math.floor(Math.random() * ITEM_COLORS.length);
  const geo = new THREE.SphereGeometry(0.35, 32, 32);
  const mat = new THREE.MeshStandardMaterial({
    color: ITEM_COLORS[i],
    emissive: ITEM_EMISSIVES[i],
    emissiveIntensity: 0.7,
    metalness: 0.9,
    roughness: 0.15,
  });
  const item = new THREE.Mesh(geo, mat);
  item.position.copy(randomItemPos());
  item.castShadow = true;
  scene.add(item);
  pickupItemList.push(item);
  itemCount.value = pickupItemList.length;
}

function createPickupItems() {
  for (let i = 0; i < 4; i++) spawnOneItem(i);
}

function removePickupItem(item: THREE.Mesh) {
  const idx = pickupItemList.indexOf(item);
  if (idx > -1) pickupItemList.splice(idx, 1);
  scene.remove(item);
  item.geometry.dispose();
  item.material.dispose();
  itemCount.value = pickupItemList.length;
}

function scheduleRespawn() {
  const timer = window.setTimeout(() => {
    const i = respawnTimers.indexOf(timer);
    if (i > -1) respawnTimers.splice(i, 1);
    if (!disposed) spawnOneItem();
  }, 3500);
  respawnTimers.push(timer);
}

// ========== 障碍物 ==========
function createObstacles() {
  const configs = [
    { x: -6, z: -5, r: 0.9, h: 2.2 },
    { x: 6, z: -6, r: 0.7, h: 1.6 },
    { x: -5, z: 6, r: 0.8, h: 1.9 },
    { x: 7, z: 5, r: 1.0, h: 2.6 },
    { x: 0, z: -8, r: 0.6, h: 1.3 },
  ];
  configs.forEach((cfg) => {
    const mat = new THREE.MeshStandardMaterial({
      color: 0x6fa8dc,
      emissive: 0x2a6db5,
      emissiveIntensity: 0.25,
      metalness: 0.6,
      roughness: 0.35,
    });
    const mesh = new THREE.Mesh(new THREE.CylinderGeometry(cfg.r, cfg.r, cfg.h, 24), mat);
    mesh.position.set(cfg.x, cfg.h / 2, cfg.z);
    mesh.castShadow = true;
    mesh.receiveShadow = true;
    scene.add(mesh);
    obstacleMeshes.push(mesh);
    obstacleList.push({ x: cfg.x, z: cfg.z, r: cfg.r });
  });
}

function resolveCollisions(model: THREE.Group) {
  for (const ob of obstacleList) {
    const dx = model.position.x - ob.x;
    const dz = model.position.z - ob.z;
    const dist = Math.hypot(dx, dz);
    const minD = ob.r + 0.7;
    if (dist < minD && dist > 0.0001) {
      model.position.x = ob.x + (dx / dist) * minD;
      model.position.z = ob.z + (dz / dist) * minD;
    }
  }
  model.position.x = THREE.MathUtils.clamp(model.position.x, -BOUND, BOUND);
  model.position.z = THREE.MathUtils.clamp(model.position.z, -BOUND, BOUND);
}

// ========== 头顶临时气泡 ==========
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
  label.position.set(0, 3.0, 0);
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

// ========== 持久名牌 ==========
function createNameTag(model: THREE.Group, name: string) {
  const div = document.createElement("div");
  div.className = "bot-tag";
  div.textContent = name;
  const label = new CSS2DObject(div);
  label.position.set(0, 2.45, 0);
  model.add(label);
  nameTagMap.set(model, label);
}

function updateTagHighlight() {
  nameTagMap.forEach((tag, model) => {
    const el = tag.element as HTMLDivElement;
    el.classList.toggle("selected", model === selectedModel);
  });
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

function selectModel(model: THREE.Group) {
  if (selectedModel === model) return;
  selectedModel = model;
  if (outlinePass) outlinePass.selectedObjects = [model];
  moveTargetMap.delete(model);
  removeRippleByModel(model);
  updateTagHighlight();
}

// ========== 射线拾取 ==========
const raycaster = new THREE.Raycaster();
const mouse = new THREE.Vector2();
const dragPlane = new THREE.Plane(new THREE.Vector3(0, 1, 0), 0);
let isDragging = false;
const tempWorldPoint = new THREE.Vector3();

function onMouseDown(e: MouseEvent) {
  if (!renderer?.domElement || modelList.length === 0) return;
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
  if (!renderer?.domElement) return;
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
      updateTagHighlight();
      dragOffset.x = selectedModel.position.x - downHitPoint!.x;
      dragOffset.y = selectedModel.position.z - downHitPoint!.z;
    }
  }

  if (isDragging && selectedModel) {
    raycaster.setFromCamera(mouse, camera);
    const intersectPoint = new THREE.Vector3();
    raycaster.ray.intersectPlane(dragPlane, intersectPoint);
    if (intersectPoint) {
      selectedModel.position.x = THREE.MathUtils.clamp(intersectPoint.x + dragOffset.x, -BOUND, BOUND);
      selectedModel.position.z = THREE.MathUtils.clamp(intersectPoint.z + dragOffset.y, -BOUND, BOUND);
      if (selectedModel.position.y < 0) selectedModel.position.y = 0;
    }
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
        if (!selectedModel) selectModel(targetModel);
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
        selectModel(downHitModel);
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
        if (!selectedModel) selectModel(targetModel);
        createModelLabel(targetModel, "我来啦！！！");
        createClickRipple(downGroundPoint.x, downGroundPoint.z, targetModel);
        moveTargetMap.set(targetModel, {
          x: THREE.MathUtils.clamp(downGroundPoint.x, -BOUND, BOUND),
          z: THREE.MathUtils.clamp(downGroundPoint.z, -BOUND, BOUND),
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

// ========== 键盘 ==========
function onKeyDown(e: KeyboardEvent) {
  const k = e.key.toLowerCase();
  if (KEY_MAP[k]) {
    keys.add(k);
    e.preventDefault();
  }
}
function onKeyUp(e: KeyboardEvent) {
  keys.delete(e.key.toLowerCase());
}

// ========== 触屏事件兼容 ==========
function onTouchStart(e: TouchEvent) {
  e.preventDefault();
  const t = e.touches[0];
  onMouseDown(new MouseEvent("mousedown", { clientX: t.clientX, clientY: t.clientY }));
}
function onTouchMove(e: TouchEvent) {
  e.preventDefault();
  const t = e.touches[0];
  onMouseMove(new MouseEvent("mousemove", { clientX: t.clientX, clientY: t.clientY }));
}
function onTouchEnd(e: TouchEvent) {
  e.preventDefault();
  onMouseUp(new MouseEvent("mouseup", { clientX: 0, clientY: 0 }));
}

// ========== 昼夜切换 ==========
function toggleNight() {
  isNight.value = !isNight.value;
  const theme = isNight.value ? NIGHT : DAY;
  scene.background = new THREE.Color(theme.bg);
  scene.fog = new THREE.FogExp2(theme.fog, theme.fogDensity);
  dirLight.color.set(theme.dirColor);
  dirLight.intensity = theme.dirIntensity;
  hemiLight.intensity = theme.hemiIntensity;
  rimLight.color.set(theme.rimColor);
  bloomPass.strength = theme.bloom;
}

// ========== 初始化场景 ==========
function init() {
  if (!containerRef.value) return;
  disposed = false;
  clock.getDelta(); // 丢弃与上次挂载之间的巨大间隔

  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;

  // 重置跨挂载残留状态
  selectedModel = null;
  modelWrapList = [];
  modelList = [];
  pickupItemList = [];
  obstacleList = [];
  obstacleMeshes = [];
  rippleList = [];
  modelAnimQueue = [];
  moveTargetMap.clear();
  nameTagMap.clear();
  keys.clear();

  scene = new THREE.Scene();
  scene.background = new THREE.Color(DAY.bg);
  scene.fog = new THREE.FogExp2(DAY.fog, DAY.fogDensity);

  camera = new THREE.PerspectiveCamera(75, w / h, 0.1, 2000);
  camera.position.set(12, 12, 12);

  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(w, h);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  renderer.toneMapping = THREE.ACESFilmicToneMapping;
  renderer.toneMappingExposure = 1.0;
  containerRef.value.appendChild(renderer.domElement);
  renderer.domElement.style.touchAction = "none";
  renderer.domElement.style.userSelect = "none";
  renderer.domElement.style.webkitTapHighlightColor = "transparent";

  labelRenderer = new CSS2DRenderer();
  labelRenderer.setSize(w, h);
  labelRenderer.domElement.style.position = "absolute";
  labelRenderer.domElement.style.top = "0";
  labelRenderer.domElement.style.left = "0";
  labelRenderer.domElement.style.pointerEvents = "none";
  containerRef.value.appendChild(labelRenderer.domElement);

  // 灯光
  dirLight = new THREE.DirectionalLight(DAY.dirColor, DAY.dirIntensity);
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

  hemiLight = new THREE.HemisphereLight(0xb8dcff, 0xd7ecff, DAY.hemiIntensity);
  scene.add(hemiLight);

  rimLight = new THREE.DirectionalLight(DAY.rimColor, 0.4);
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

  // 边界提示框
  const boundGeo = new THREE.PlaneGeometry(BOUND * 2, BOUND * 2);
  const boundMat = new THREE.MeshBasicMaterial({
    color: 0x4da6ff,
    transparent: true,
    opacity: 0.06,
    side: THREE.DoubleSide,
    depthWrite: false,
  });
  const boundPlane = new THREE.Mesh(boundGeo, boundMat);
  boundPlane.rotation.x = -Math.PI / 2;
  boundPlane.position.y = 0.03;
  scene.add(boundPlane);

  // 轨道控制器
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.05;
  controls.minDistance = 2;
  controls.maxDistance = 60;
  controls.maxPolarAngle = Math.PI * 0.49;

  // 鼠标事件
  renderer.domElement.addEventListener("mousedown", onMouseDown);
  renderer.domElement.addEventListener("mousemove", onMouseMove);
  renderer.domElement.addEventListener("mouseup", onMouseUp);
  renderer.domElement.addEventListener("mouseleave", onMouseUp);

  // 触屏事件
  renderer.domElement.addEventListener("touchstart", onTouchStart, { passive: false });
  renderer.domElement.addEventListener("touchmove", onTouchMove, { passive: false });
  renderer.domElement.addEventListener("touchend", onTouchEnd, { passive: false });

  // 键盘
  window.addEventListener("keydown", onKeyDown);
  window.addEventListener("keyup", onKeyUp);

  // ========== 后期处理管线 ==========
  composer = new EffectComposer(renderer);
  composer.addPass(new RenderPass(scene, camera));

  ssaoPass = new SSAOPass(scene, camera, w, h);
  ssaoPass.kernelRadius = 8;
  ssaoPass.minDistance = 0.005;
  ssaoPass.maxDistance = 0.05;
  ssaoPass.output = SSAOPass.OUTPUT.Default;
  composer.addPass(ssaoPass);

  outlinePass = new OutlinePass(new THREE.Vector2(w, h), scene, camera);
  outlinePass.edgeStrength = 3.0;
  outlinePass.edgeGlow = 0.65;
  outlinePass.edgeThickness = 1.2;
  outlinePass.pulsePeriod = 0;
  outlinePass.visibleEdgeColor.set(0x4da6ff);
  outlinePass.hiddenEdgeColor.set(0x000000);
  composer.addPass(outlinePass);

  bloomPass = new UnrealBloomPass(new THREE.Vector2(w, h), DAY.bloom, 0.5, 0.85);
  composer.addPass(bloomPass);

  composer.addPass(new OutputPass());

  // 场景内容
  createPickupItems();
  createObstacles();
  initParticlePool();

  // ========== 加载 3 个角色 ==========
  loadingCount.value = 3;
  const charConfigs = [
    { pos: new THREE.Vector3(-4, 0, 1), tint: 0xffffff, name: "小白" },
    { pos: new THREE.Vector3(0, 0, 2.5), tint: 0x9fc8ff, name: "小蓝" },
    { pos: new THREE.Vector3(4, 0, 1), tint: 0xffb3c8, name: "小粉" },
  ];
  charConfigs.forEach((cfg) => loadModelWithAnim("/scene.glb", cfg.pos, cfg.tint, cfg.name));

  // 尺寸监听（tab 切换 / 响应式均可靠）
  resizeObserver = new ResizeObserver(() => onResize());
  resizeObserver.observe(containerRef.value);

  animate();
}

function loadModelWithAnim(url: string, pos: THREE.Vector3, tint: number, name: string) {
  const loader = new GLTFLoader();
  loader.load(
    url,
    (gltf) => {
      if (disposed || !scene) return; // 组件可能已卸载
      const model = gltf.scene;
      model.position.copy(pos);
      if (model.position.y < 0) model.position.y = 0;
      model.scale.set(BASE_SCALE, BASE_SCALE, BASE_SCALE);
      model.traverse((obj) => {
        obj.castShadow = true;
        obj.receiveShadow = true;
        const m = obj as THREE.Mesh;
        if (m.isMesh && m.material && tint !== 0xffffff) {
          // 克隆材质独立着色，避免多个角色共享材质互相影响
          if (Array.isArray(m.material)) {
            m.material = m.material.map((s) => {
              const c = s.clone();
              c.color = new THREE.Color(tint);
              return c;
            });
          } else {
            const c = m.material.clone();
            c.color = new THREE.Color(tint);
            m.material = c;
          }
        }
      });
      scene.add(model);
      modelList.push(model);
      createNameTag(model, name);

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

      loadingCount.value = Math.max(0, loadingCount.value - 1);

      // 默认选中第一个加载完成的角色
      if (!selectedModel) {
        selectModel(model);
        createModelLabel(model, "选中了我！");
      }
    },
    (p) => {
      if (p.total) {
        // 可选：控制台观察进度
      }
    },
    (err) => {
      console.error("模型加载失败", err);
      loadingCount.value = Math.max(0, loadingCount.value - 1);
    }
  );
}

function onResize() {
  if (!containerRef.value || !renderer) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;
  if (w === 0 || h === 0) return;
  camera.aspect = w / h;
  camera.updateProjectionMatrix();
  renderer.setSize(w, h);
  composer.setSize(w, h);
  labelRenderer.setSize(w, h);
  if (ssaoPass) ssaoPass.setSize(w, h);
  if (bloomPass) bloomPass.setSize(w, h);
}

function animate() {
  animationId = requestAnimationFrame(animate);
  controls.update();
  const delta = Math.min(clock.getDelta(), 0.05); // 防切页后巨大 delta
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

  // 5. WASD 键盘控制选中角色
  if (selectedModel && !isDragMode && modelList.includes(selectedModel)) {
    let kx = 0;
    let kz = 0;
    keys.forEach((k) => {
      const d = KEY_MAP[k];
      if (d) {
        kx += d[0];
        kz += d[1];
      }
    });
    if (kx !== 0 || kz !== 0) {
      moveTargetMap.delete(selectedModel);
      removeRippleByModel(selectedModel);
      const len = Math.hypot(kx, kz);
      const spd = MOVE_SPEED * 1.35 * delta;
      selectedModel.position.x += (kx / len) * spd;
      selectedModel.position.z += (kz / len) * spd;
      selectedModel.rotation.y = Math.atan2(kx, kz);
      playRun(selectedModel, true);
      keyboardMoving = true;
    } else if (keyboardMoving) {
      keyboardMoving = false;
      playRun(selectedModel, false);
    }
  }

  // 6. NPC 自动巡逻
  if (wanderAI.value) {
    modelList.forEach((m) => {
      if (m === selectedModel) return;
      if (!moveTargetMap.has(m) && Math.random() < delta * 0.3) {
        moveTargetMap.set(m, {
          x: (Math.random() * 2 - 1) * (BOUND - 1),
          z: (Math.random() * 2 - 1) * (BOUND - 1),
        });
        playRun(m, true);
      }
    });
  }

  // 7. 自动移动（点击目标 / 巡逻目标共用）
  moveTargetMap.forEach((target, model) => {
    const dx = target.x - model.position.x;
    const dz = target.z - model.position.z;
    const distance = Math.sqrt(dx * dx + dz * dz);
    if (distance < 0.05) {
      model.position.x = target.x;
      model.position.z = target.z;
      moveTargetMap.delete(model);
      playRun(model, false);

      if (target.pickupItem && pickupItemList.includes(target.pickupItem)) {
        removePickupItem(target.pickupItem);
        score.value += 1;
        createModelLabel(model, `+1 得分 ${score.value}`);
        spawnExplode(model.position.clone().setY(0.8), 0.8);
        scheduleRespawn();
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

  // 8. 障碍碰撞 + 边界（所有角色）
  modelList.forEach((model) => resolveCollisions(model));

  // 9. 弹跳动画
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

  // 10. 光圈特效
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

  // 11. 跟随相机
  if (followCam.value && selectedModel) {
    const t = new THREE.Vector3(
      selectedModel.position.x,
      selectedModel.position.y + 1,
      selectedModel.position.z
    );
    controls.target.lerp(t, 0.08);
  }

  // 12. Y 轴钳制
  modelList.forEach((model) => {
    if (model.position.y < 0) model.position.y = 0;
  });

  // 渲染
  composer.render();
  labelRenderer.render(scene, camera);
}

function dispose() {
  disposed = true;
  if (animationId !== null) cancelAnimationFrame(animationId);
  resizeObserver?.disconnect();
  window.removeEventListener("keydown", onKeyDown);
  window.removeEventListener("keyup", onKeyUp);

  if (renderer?.domElement) {
    renderer.domElement.removeEventListener("mousedown", onMouseDown);
    renderer.domElement.removeEventListener("mousemove", onMouseMove);
    renderer.domElement.removeEventListener("mouseup", onMouseUp);
    renderer.domElement.removeEventListener("mouseleave", onMouseUp);
    renderer.domElement.removeEventListener("touchstart", onTouchStart);
    renderer.domElement.removeEventListener("touchmove", onTouchMove);
    renderer.domElement.removeEventListener("touchend", onTouchEnd);
  }

  controls?.dispose();
  composer?.dispose();
  renderer?.dispose();

  disposeParticlePool();

  respawnTimers.forEach((t) => clearTimeout(t));
  respawnTimers.length = 0;

  labelTimeoutMap.forEach((timer) => clearTimeout(timer));
  labelTimeoutMap.clear();
  labelMap.forEach((label, model) => model.remove(label));
  labelMap.clear();
  nameTagMap.forEach((label, model) => model.remove(label));
  nameTagMap.clear();
  if (labelRenderer?.domElement) {
    labelRenderer.domElement.remove();
  }

  pickupItemList.forEach((item) => {
    scene.remove(item);
    item.geometry.dispose();
    item.material.dispose();
  });
  pickupItemList = [];

  obstacleMeshes.forEach((mesh) => {
    scene.remove(mesh);
    mesh.geometry.dispose();
    mesh.material.dispose();
  });
  obstacleMeshes = [];
  obstacleList = [];

  rippleList.forEach((item) => {
    scene.remove(item.mesh);
    item.mesh.geometry.dispose();
    item.mesh.material.dispose();
  });
  rippleList = [];

  modelAnimQueue = [];
  moveTargetMap.clear();
  modelWrapList.forEach((wrap) => wrap.mixer.stopAllAction());
  modelWrapList = [];
  selectedModel = null;

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
  modelList = [];
}

onMounted(() => init());
onUnmounted(() => dispose());
</script>

<style scoped>
.bot-demo {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  background: #e6f4ff;
  touch-action: none;
  user-select: none;
}

.bot-toolbar {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  display: flex;
  gap: 8px;
}

.bot-toolbar button {
  padding: 5px 12px;
  font-size: 0.75rem;
  border: 1px solid rgba(77, 166, 255, 0.4);
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.75);
  color: #2a6db5;
  cursor: pointer;
  backdrop-filter: blur(6px);
  transition: all 0.2s;
}

.bot-toolbar button:hover {
  border-color: #4da6ff;
  color: #14406e;
}

.bot-toolbar button.active {
  background: linear-gradient(135deg, #4da6ff, #7c6cff);
  border-color: #4da6ff;
  color: #fff;
}

.bot-hud {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 10;
  display: flex;
  gap: 12px;
  padding: 6px 14px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.75);
  border: 1px solid rgba(77, 166, 255, 0.35);
  backdrop-filter: blur(6px);
  font-size: 0.78rem;
  color: #2a6db5;
}

.bot-hud b {
  color: #14406e;
  font-size: 0.9rem;
}

.bot-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(230, 244, 255, 0.85);
  color: #2a6db5;
  font-size: 1rem;
  z-index: 20;
  letter-spacing: 0.05em;
}

/* 持久名牌（CSS2DRenderer 渲染到容器内，需要全局样式） */
:global(.bot-tag) {
  padding: 3px 10px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(77, 166, 255, 0.5);
  color: #2a6db5;
  font-size: 12px;
  font-weight: 600;
  font-family: system-ui, sans-serif;
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 2px 8px rgba(77, 166, 255, 0.25);
}

:global(.bot-tag.selected) {
  background: linear-gradient(135deg, #4da6ff, #7c6cff);
  border-color: #4da6ff;
  color: #fff;
  box-shadow: 0 0 14px rgba(124, 108, 255, 0.55);
}
</style>
