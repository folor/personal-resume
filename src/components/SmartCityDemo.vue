<template>
  <!-- 智慧城市 3D 场景容器 -->
  <div ref="containerRef" class="city-container">
    <div class="city-toolbar">
      <button :class="{ active: autoRotate }" @click="toggleRotate">自动巡航</button>
      <button :class="{ active: dataFlowOn }" @click="toggleFlow">数据流</button>
      <button @click="resetView">重置视角</button>
    </div>
  </div>
</template>

<script setup lang="ts">
// @ts-nocheck
import { ref, onMounted, onUnmounted } from "vue";
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";
import { EffectComposer } from "three/examples/jsm/postprocessing/EffectComposer.js";
import { RenderPass } from "three/examples/jsm/postprocessing/RenderPass.js";
import { UnrealBloomPass } from "three/examples/jsm/postprocessing/UnrealBloomPass.js";
import { OutputPass } from "three/examples/jsm/postprocessing/OutputPass.js";
import { CSS2DRenderer, CSS2DObject } from "three/examples/jsm/renderers/CSS2DRenderer.js";

const containerRef = ref<HTMLDivElement | null>(null);
const autoRotate = ref(true);
const dataFlowOn = ref(true);

let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let labelRenderer: CSS2DRenderer;
let controls: OrbitControls;
let composer: EffectComposer;
let bloomPass: UnrealBloomPass;
let animationId: number | null = null;
let resizeObserver: ResizeObserver | null = null;
const clock = new THREE.Clock();

// 城市元素
const buildingList: THREE.Mesh[] = [];
const buildingInfoMap = new Map<THREE.Mesh, BuildingInfo>();
const carList: { mesh: THREE.Mesh; axis: "x" | "z"; fixed: number; pos: number; speed: number; dir: number }[] = [];
const beaconMats: THREE.MeshStandardMaterial[] = [];

// 数据流
interface FlowItem {
  curve: THREE.QuadraticBezierCurve3;
  t: number;
  speed: number;
  pointIndex: number;
}
let flowPoints: THREE.Points | null = null;
let flowItems: FlowItem[] = [];
const FLOW_COUNT = 46;

// 交互
const raycaster = new THREE.Raycaster();
const pointer = new THREE.Vector2();
let hoveredBuilding: THREE.Mesh | null = null;
let hoverLabel: CSS2DObject | null = null;
const matSwapCache = new Map<THREE.Material, THREE.Material>();
let originalMatMap = new Map<THREE.Mesh, THREE.Material>();
let focusTween: { camPos: THREE.Vector3; target: THREE.Vector3 } | null = null;
const DEFAULT_CAM = new THREE.Vector3(95, 75, 95);
const DEFAULT_TARGET = new THREE.Vector3(0, 8, 0);
let downX = 0;
let downY = 0;

interface BuildingInfo {
  name: string;
  floors: number;
  power: number;
  occupancy: number;
  traffic: number;
  type: string;
}

const BUILDING_TYPES = ["商务写字楼", "智慧数据中心", "科技园区", "商业综合体", "智慧楼宇", "政务中心"];

// ========== 生成窗户发光贴图 ==========
function makeWindowTexture(hueBase: number): THREE.Texture {
  const canvas = document.createElement("canvas");
  canvas.width = 64;
  canvas.height = 128;
  const ctx = canvas.getContext("2d")!;
  ctx.fillStyle = "#04070f";
  ctx.fillRect(0, 0, 64, 128);
  const cols = 6;
  const rows = 16;
  const cw = 64 / cols;
  const ch = 128 / rows;
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      if (Math.random() < 0.42) continue; // 部分窗户熄灭
      const hue = hueBase + (Math.random() - 0.5) * 30;
      const light = 0.55 + Math.random() * 0.45;
      ctx.fillStyle = `hsl(${hue}, 85%, ${light * 60}%)`;
      ctx.fillRect(c * cw + 1.5, r * ch + 1.5, cw - 3, ch - 3);
    }
  }
  const tex = new THREE.CanvasTexture(canvas);
  tex.colorSpace = THREE.SRGBColorSpace;
  return tex;
}

// ========== 生成城市 ==========
function buildCity() {
  // 地面
  const ground = new THREE.Mesh(
    new THREE.PlaneGeometry(320, 320),
    new THREE.MeshStandardMaterial({ color: 0x070d1a, roughness: 0.9, metalness: 0.2 })
  );
  ground.rotation.x = -Math.PI / 2;
  ground.receiveShadow = true;
  scene.add(ground);

  // 网格
  const grid = new THREE.GridHelper(320, 64, 0x14365c, 0x0d2038);
  grid.position.y = 0.02;
  scene.add(grid);

  // 道路
  const roadMat = new THREE.MeshStandardMaterial({ color: 0x101a2b, roughness: 0.6, metalness: 0.3 });
  const roadPositions = [-60, -40, -20, 0, 20, 40, 60];
  roadPositions.forEach((p) => {
    const rx = new THREE.Mesh(new THREE.PlaneGeometry(320, 8), roadMat);
    rx.rotation.x = -Math.PI / 2;
    rx.position.set(0, 0.03, p);
    scene.add(rx);
    const rz = new THREE.Mesh(new THREE.PlaneGeometry(8, 320), roadMat);
    rz.rotation.x = -Math.PI / 2;
    rz.position.set(p, 0.03, 0);
    scene.add(rz);
  });

  // 楼体材质（几种贴图变体共享）
  const windowTextures = [
    makeWindowTexture(200),
    makeWindowTexture(190),
    makeWindowTexture(210),
    makeWindowTexture(185),
  ];
  const buildingMats = windowTextures.map(
    (t) =>
      new THREE.MeshStandardMaterial({
        color: 0x0b1626,
        emissive: 0xffffff,
        emissiveMap: t,
        emissiveIntensity: 0.9,
        roughness: 0.4,
        metalness: 0.6,
        map: t,
      })
  );

  const blockCenters = [-50, -30, -10, 10, 30, 50];
  let landmarkCount = 0;

  for (const bx of blockCenters) {
    for (const bz of blockCenters) {
      if (bx === 10 && bz === 10) continue; // 中心留给数据塔（近似中心块）

      const distToCenter = Math.sqrt(bx * bx + bz * bz);
      const downtownFactor = Math.max(0, 1 - distToCenter / 90); // 越靠中心越高

      // 少量地标建筑
      const isLandmark = landmarkCount < 4 && downtownFactor > 0.55 && Math.random() < 0.25;
      if (isLandmark) landmarkCount++;

      if (isLandmark) {
        const h = 52 + Math.random() * 22;
        addBuilding(bx, bz, 8 + Math.random() * 3, h, 8 + Math.random() * 3, buildingMats, true);
      } else {
        // 每个街区 2~4 栋楼
        const count = 2 + Math.floor(Math.random() * 3);
        const offsets = [-4.5, 0, 4.5];
        for (let i = 0; i < count; i++) {
          const ox = offsets[Math.floor(Math.random() * offsets.length)];
          const oz = offsets[Math.floor(Math.random() * offsets.length)];
          const h = 8 + downtownFactor * 32 * (0.5 + Math.random() * 0.8) + Math.random() * 8;
          const w = 3.5 + Math.random() * 3;
          const d = 3.5 + Math.random() * 3;
          addBuilding(bx + ox, bz + oz, w, h, d, buildingMats, h > 38);
        }
      }
    }
  }

  buildCentralTower();
}

function addBuilding(
  x: number,
  z: number,
  w: number,
  h: number,
  d: number,
  mats: THREE.MeshStandardMaterial[],
  withBeacon: boolean
) {
  const geo = new THREE.BoxGeometry(w, h, d);
  const mat = mats[Math.floor(Math.random() * mats.length)];
  const mesh = new THREE.Mesh(geo, mat);
  mesh.position.set(x, h / 2, z);
  mesh.castShadow = true;
  mesh.receiveShadow = true;
  scene.add(mesh);
  buildingList.push(mesh);

  const floors = Math.max(2, Math.floor(h / 3));
  buildingInfoMap.set(mesh, {
    name: `${String.fromCharCode(65 + Math.floor(Math.random() * 26))}-${Math.floor(Math.random() * 90 + 10)} 号楼`,
    floors,
    power: Math.floor(h * (18 + Math.random() * 20)),
    occupancy: Math.floor(62 + Math.random() * 36),
    traffic: Math.floor(h * (30 + Math.random() * 60)),
    type: BUILDING_TYPES[Math.floor(Math.random() * BUILDING_TYPES.length)],
  });

  // 顶部警示灯
  if (withBeacon) {
    const beaconMat = new THREE.MeshStandardMaterial({
      color: 0xff3344,
      emissive: 0xff2233,
      emissiveIntensity: 2,
    });
    beaconMats.push(beaconMat);
    const beacon = new THREE.Mesh(new THREE.SphereGeometry(0.35, 12, 12), beaconMat);
    beacon.position.set(x, h + 0.6, z);
    scene.add(beacon);
    const rod = new THREE.Mesh(
      new THREE.CylinderGeometry(0.06, 0.06, 1.2, 6),
      new THREE.MeshStandardMaterial({ color: 0x334455 })
    );
    rod.position.set(x, h + 0.1, z);
    scene.add(rod);
  }
}

// ========== 中央数据塔 ==========
let towerGroup: THREE.Group;
function buildCentralTower() {
  towerGroup = new THREE.Group();
  towerGroup.position.set(0, 0, 0);

  const towerMat = new THREE.MeshStandardMaterial({
    color: 0x0d2038,
    emissive: 0x1899ff,
    emissiveIntensity: 0.35,
    metalness: 0.85,
    roughness: 0.25,
  });
  const ringMat = new THREE.MeshStandardMaterial({
    color: 0x0a1830,
    emissive: 0x33ccff,
    emissiveIntensity: 2.2,
  });

  const core = new THREE.Mesh(new THREE.CylinderGeometry(2.2, 3.4, 72, 8), towerMat);
  core.position.y = 36;
  towerGroup.add(core);

  for (let i = 0; i < 7; i++) {
    const r = 4.6 - i * 0.45;
    const ring = new THREE.Mesh(new THREE.TorusGeometry(r, 0.16, 8, 40), ringMat);
    ring.rotation.x = Math.PI / 2;
    ring.position.y = 10 + i * 9;
    towerGroup.add(ring);
  }

  const top = new THREE.Mesh(new THREE.SphereGeometry(1.4, 16, 16), ringMat);
  top.position.y = 73.5;
  towerGroup.add(top);

  const beamMat = new THREE.MeshBasicMaterial({
    color: 0x33ccff,
    transparent: true,
    opacity: 0.16,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
  });
  const beam = new THREE.Mesh(new THREE.CylinderGeometry(0.8, 3.2, 130, 16, 1, true), beamMat);
  beam.position.y = 135;
  towerGroup.add(beam);

  scene.add(towerGroup);
}

// ========== 车流 ==========
function buildTraffic() {
  const carGeoX = new THREE.BoxGeometry(2.0, 0.4, 0.8);
  const roads = [-60, -40, -20, 0, 20, 40, 60];

  for (let i = 0; i < 64; i++) {
    const dir = Math.random() < 0.5 ? 1 : -1;
    const axis: "x" | "z" = Math.random() < 0.5 ? "x" : "z";
    const fixed = roads[Math.floor(Math.random() * roads.length)] + dir * 1.8;
    // 车头灯（暖白）/ 尾灯（红）——对面车流
    const mat = new THREE.MeshStandardMaterial({
      color: dir > 0 ? 0xfff2cc : 0xff5533,
      emissive: dir > 0 ? 0xffeebb : 0xff3311,
      emissiveIntensity: 2.4,
    });
    const mesh = new THREE.Mesh(carGeoX, mat);
    if (axis === "z") mesh.rotation.y = Math.PI / 2;
    mesh.position.y = 0.5;
    scene.add(mesh);
    carList.push({
      mesh,
      axis,
      fixed,
      pos: (Math.random() * 2 - 1) * 150,
      speed: 22 + Math.random() * 26,
      dir,
    });
  }
}

function tickTraffic(delta: number) {
  for (const car of carList) {
    car.pos += car.speed * car.dir * delta;
    if (car.pos > 160) car.pos = -160;
    if (car.pos < -160) car.pos = 160;
    if (car.axis === "x") car.mesh.position.set(car.pos, 0.5, car.fixed);
    else car.mesh.position.set(car.fixed, 0.5, car.pos);
  }
}

// ========== 数据流（建筑 → 中央塔） ==========
function buildDataFlow() {
  const tex = new THREE.PointsMaterial({
    color: 0x44ddff,
    size: 0.9,
    transparent: true,
    opacity: 0.95,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
    sizeAttenuation: true,
  });
  const posArr = new Float32Array(FLOW_COUNT * 3);
  const geo = new THREE.BufferGeometry();
  geo.setAttribute("position", new THREE.BufferAttribute(posArr, 3));
  flowPoints = new THREE.Points(geo, tex);
  scene.add(flowPoints);

  const lineMat = new THREE.LineBasicMaterial({
    color: 0x2299cc,
    transparent: true,
    opacity: 0.12,
    blending: THREE.AdditiveBlending,
    depthWrite: false,
  });

  const sources = buildingList.filter(() => Math.random() < 0.5);
  for (let i = 0; i < FLOW_COUNT; i++) {
    const src = sources.length ? sources[i % sources.length] : buildingList[i % buildingList.length];
    const start = new THREE.Vector3(src.position.x, src.position.y + 2, src.position.z);
    const end = new THREE.Vector3(0, 68, 0);
    const mid = start.clone().lerp(end, 0.5);
    mid.y = Math.max(start.y, end.y) * 0.55 + 18;
    mid.x += (Math.random() - 0.5) * 24;
    mid.z += (Math.random() - 0.5) * 24;
    const curve = new THREE.QuadraticBezierCurve3(start, mid, end);

    const line = new THREE.Line(new THREE.BufferGeometry().setFromPoints(curve.getPoints(36)), lineMat);
    scene.add(line);

    flowItems.push({
      curve,
      t: Math.random(),
      speed: 0.12 + Math.random() * 0.22,
      pointIndex: i,
    });
  }
}

function tickDataFlow(delta: number) {
  if (!flowPoints) return;
  const attr = flowPoints.geometry.attributes.position as THREE.BufferAttribute;
  if (!dataFlowOn.value) {
    // 数据流关闭时粒子渐隐（移到塔顶汇合）
    for (const it of flowItems) {
      it.t = Math.min(1, it.t + delta * 0.5);
      const p = it.curve.getPoint(it.t);
      attr.setXYZ(it.pointIndex, p.x, p.y, p.z);
    }
    attr.needsUpdate = true;
    return;
  }
  for (const it of flowItems) {
    it.t += it.speed * delta;
    if (it.t > 1) it.t = 0;
    const p = it.curve.getPoint(it.t);
    attr.setXYZ(it.pointIndex, p.x, p.y, p.z);
  }
  attr.needsUpdate = true;
}

// ========== 交互：hover 高亮 + 信息牌 ==========
function getHighlightMat(origin: THREE.Material): THREE.Material {
  if (!matSwapCache.has(origin)) {
    const clone = (origin as THREE.MeshStandardMaterial).clone();
    clone.emissive = new THREE.Color(0x33bbff);
    clone.emissiveIntensity = 1.6;
    clone.color = new THREE.Color(0x1a4a7a);
    matSwapCache.set(origin, clone);
  }
  return matSwapCache.get(origin)!;
}

function clearHover() {
  if (hoveredBuilding && originalMatMap.has(hoveredBuilding)) {
    hoveredBuilding.material = originalMatMap.get(hoveredBuilding)!;
  }
  hoveredBuilding = null;
  if (hoverLabel) {
    scene.remove(hoverLabel);
    hoverLabel = null;
  }
  if (containerRef.value) containerRef.value.style.cursor = "grab";
}

function applyHover(mesh: THREE.Mesh) {
  hoveredBuilding = mesh;
  if (!originalMatMap.has(mesh)) originalMatMap.set(mesh, mesh.material);
  mesh.material = getHighlightMat(mesh.material);
  if (containerRef.value) containerRef.value.style.cursor = "pointer";

  const info = buildingInfoMap.get(mesh);
  const div = document.createElement("div");
  div.className = "city-label";
  div.innerHTML = `
    <div class="city-label-title">${info.name} <span>${info.type}</span></div>
    <div class="city-label-row"><i>楼层</i><b>${info.floors} 层</b></div>
    <div class="city-label-row"><i>实时功耗</i><b>${info.power} kW</b></div>
    <div class="city-label-row"><i>入驻率</i><b>${info.occupancy}%</b></div>
    <div class="city-label-row"><i>今日人流</i><b>${info.traffic.toLocaleString()}</b></div>
  `;
  hoverLabel = new CSS2DObject(div);
  hoverLabel.position.set(mesh.position.x, mesh.geometry.parameters.height + 3, mesh.position.z);
  scene.add(hoverLabel);
}

function updatePointer(e: PointerEvent) {
  if (!renderer?.domElement) return;
  const rect = renderer.domElement.getBoundingClientRect();
  pointer.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
  pointer.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
}

function onPointerMove(e: PointerEvent) {
  updatePointer(e);
  if (buildingList.length === 0) return;
  raycaster.setFromCamera(pointer, camera);
  const hits = raycaster.intersectObjects(buildingList, false);
  const first = hits.length > 0 ? (hits[0].object as THREE.Mesh) : null;

  if (first === hoveredBuilding) return;
  clearHover();
  if (first) applyHover(first);
}

function onPointerDown(e: PointerEvent) {
  downX = e.clientX;
  downY = e.clientY;
}

function onPointerUp(e: PointerEvent) {
  // 区分点击与拖拽
  const dx = Math.abs(e.clientX - downX);
  const dy = Math.abs(e.clientY - downY);
  if (dx > 6 || dy > 6) return;

  updatePointer(e);
  raycaster.setFromCamera(pointer, camera);
  const hits = raycaster.intersectObjects(buildingList, false);
  if (hits.length > 0) {
    // 点击建筑 → 聚焦
    const b = hits[0].object as THREE.Mesh;
    const h = b.geometry.parameters.height;
    const dist = 22 + h * 0.4;
    const dir = new THREE.Vector3().subVectors(camera.position, controls.target).normalize();
    focusTween = {
      target: new THREE.Vector3(b.position.x, h * 0.55, b.position.z),
      camPos: new THREE.Vector3(
        b.position.x + dir.x * dist,
        h * 0.7 + dist * 0.55,
        b.position.z + dir.z * dist
      ),
    };
    controls.autoRotate = false;
    autoRotate.value = false;
  } else {
    // 点击空白 → 回默认视角
    focusTween = { target: DEFAULT_TARGET.clone(), camPos: DEFAULT_CAM.clone() };
  }
}

// ========== 工具栏 ==========
function toggleRotate() {
  autoRotate.value = !autoRotate.value;
  controls.autoRotate = autoRotate.value;
}
function toggleFlow() {
  dataFlowOn.value = !dataFlowOn.value;
}
function resetView() {
  focusTween = { target: DEFAULT_TARGET.clone(), camPos: DEFAULT_CAM.clone() };
  autoRotate.value = false;
  controls.autoRotate = false;
}

// ========== 初始化 ==========
function init() {
  if (!containerRef.value) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;

  scene = new THREE.Scene();
  scene.background = new THREE.Color(0x050a14);
  scene.fog = new THREE.FogExp2(0x050a14, 0.0035);

  camera = new THREE.PerspectiveCamera(55, w / h, 0.1, 1000);
  camera.position.copy(DEFAULT_CAM);

  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(w, h);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  renderer.toneMapping = THREE.ACESFilmicToneMapping;
  renderer.toneMappingExposure = 1.05;
  containerRef.value.appendChild(renderer.domElement);

  labelRenderer = new CSS2DRenderer();
  labelRenderer.setSize(w, h);
  labelRenderer.domElement.style.position = "absolute";
  labelRenderer.domElement.style.top = "0";
  labelRenderer.domElement.style.left = "0";
  labelRenderer.domElement.style.pointerEvents = "none";
  containerRef.value.appendChild(labelRenderer.domElement);

  // 灯光
  const moon = new THREE.DirectionalLight(0x8899ff, 0.6);
  moon.position.set(60, 100, -40);
  moon.castShadow = true;
  moon.shadow.mapSize.set(2048, 2048);
  moon.shadow.camera.left = -120;
  moon.shadow.camera.right = 120;
  moon.shadow.camera.top = 120;
  moon.shadow.camera.bottom = -120;
  scene.add(moon);
  scene.add(new THREE.HemisphereLight(0x223355, 0x05070d, 0.55));
  const cityGlow = new THREE.PointLight(0x33aaff, 1.2, 180, 1.6);
  cityGlow.position.set(0, 60, 0);
  scene.add(cityGlow);

  // 控制器
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.06;
  controls.target.copy(DEFAULT_TARGET);
  controls.minDistance = 20;
  controls.maxDistance = 320;
  controls.maxPolarAngle = Math.PI * 0.48;
  controls.autoRotate = autoRotate.value;
  controls.autoRotateSpeed = 0.5;

  // 事件
  renderer.domElement.addEventListener("pointermove", onPointerMove);
  renderer.domElement.addEventListener("pointerdown", onPointerDown);
  renderer.domElement.addEventListener("pointerup", onPointerUp);

  // 城市内容
  buildCity();
  buildTraffic();
  buildDataFlow();

  // 后期处理：Bloom 让灯光/窗户/数据流发光
  composer = new EffectComposer(renderer);
  composer.addPass(new RenderPass(scene, camera));
  bloomPass = new UnrealBloomPass(new THREE.Vector2(w, h), 0.75, 0.55, 0.18);
  composer.addPass(bloomPass);
  composer.addPass(new OutputPass());

  // 尺寸监听（tab 切换/响应式）
  resizeObserver = new ResizeObserver(() => onResize());
  resizeObserver.observe(containerRef.value);

  animate();
}

function onResize() {
  if (!containerRef.value) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;
  if (w === 0 || h === 0) return;
  camera.aspect = w / h;
  camera.updateProjectionMatrix();
  renderer.setSize(w, h);
  composer.setSize(w, h);
  labelRenderer.setSize(w, h);
}

function animate() {
  animationId = requestAnimationFrame(animate);
  const delta = Math.min(clock.getDelta(), 0.05);
  const elapsed = clock.getElapsedTime();

  tickTraffic(delta);
  tickDataFlow(delta);

  // 航空警示灯闪烁
  const blink = 1.2 + Math.sin(elapsed * 3.2) * 1.1;
  beaconMats.forEach((m) => (m.emissiveIntensity = Math.max(0.1, blink)));

  // 聚焦补间
  if (focusTween) {
    camera.position.lerp(focusTween.camPos, 0.06);
    controls.target.lerp(focusTween.target, 0.08);
    if (camera.position.distanceTo(focusTween.camPos) < 0.6) focusTween = null;
  }

  controls.update();
  composer.render();
  labelRenderer.render(scene, camera);
}

// ========== 清理 ==========
function dispose() {
  if (animationId !== null) cancelAnimationFrame(animationId);
  resizeObserver?.disconnect();
  window.removeEventListener("resize", onResize);

  if (renderer?.domElement) {
    renderer.domElement.removeEventListener("pointermove", onPointerMove);
    renderer.domElement.removeEventListener("pointerdown", onPointerDown);
    renderer.domElement.removeEventListener("pointerup", onPointerUp);
  }

  controls?.dispose();
  composer?.dispose();
  renderer?.dispose();
  if (labelRenderer?.domElement) labelRenderer.domElement.remove();

  clearHover();
  scene?.traverse((obj) => {
    const m = obj as THREE.Mesh;
    if (m.geometry) m.geometry.dispose();
    const mat = (m as any).material;
    if (mat) {
      if (Array.isArray(mat)) mat.forEach((s: any) => s.dispose());
      else {
        if (mat.map) mat.map.dispose();
        if (mat.emissiveMap) mat.emissiveMap.dispose();
        mat.dispose();
      }
    }
  });

  buildingList.length = 0;
  buildingInfoMap.clear();
  matSwapCache.clear();
  originalMatMap = new Map();
  carList.length = 0;
  flowItems = [];
  beaconMats.length = 0;
}

onMounted(() => init());
onUnmounted(() => dispose());
</script>

<style scoped>
.city-container {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: hidden;
  background: #050a14;
  cursor: grab;
  touch-action: none;
}

.city-toolbar {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 10;
  display: flex;
  gap: 8px;
}

.city-toolbar button {
  padding: 5px 12px;
  font-size: 0.75rem;
  border: 1px solid rgba(102, 204, 255, 0.35);
  border-radius: 6px;
  background: rgba(6, 16, 32, 0.72);
  color: #9fd8ff;
  cursor: pointer;
  backdrop-filter: blur(6px);
  transition: all 0.2s;
}

.city-toolbar button:hover {
  border-color: rgba(102, 204, 255, 0.8);
  color: #e6f7ff;
}

.city-toolbar button.active {
  background: rgba(30, 130, 220, 0.4);
  border-color: #4db8ff;
  color: #ffffff;
}

/* 3D 信息牌（CSS2DRenderer 挂载在组件外的 body 下的容器内，需要全局样式） */
:global(.city-label) {
  padding: 10px 14px;
  background: rgba(5, 14, 28, 0.88);
  border: 1px solid rgba(77, 184, 255, 0.55);
  border-radius: 8px;
  font-size: 12px;
  color: #cfe8ff;
  font-family: system-ui, sans-serif;
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 0 18px rgba(51, 170, 255, 0.35);
  transform: translateY(-6px);
}

:global(.city-label-title) {
  font-size: 13px;
  font-weight: 700;
  color: #7fd0ff;
  margin-bottom: 6px;
  border-bottom: 1px solid rgba(77, 184, 255, 0.25);
  padding-bottom: 4px;
}

:global(.city-label-title span) {
  font-weight: 400;
  font-size: 11px;
  color: #6f93b8;
  margin-left: 6px;
}

:global(.city-label-row) {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  line-height: 1.7;
}

:global(.city-label-row i) {
  font-style: normal;
  color: #6f93b8;
}

:global(.city-label-row b) {
  color: #e8f5ff;
  font-weight: 600;
}
</style>
