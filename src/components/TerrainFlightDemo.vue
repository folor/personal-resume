<template>
  <div ref="containerRef" class="terrain-demo">
    <!-- 加载遮罩 -->
    <div v-if="loading" class="terrain-loading">
      <div class="loading-spinner"></div>
      <p>正在生成地形地貌…</p>
    </div>

    <!-- HUD 仪表 -->
    <div class="hud" v-if="!loading">
      <div class="hud-item">
        <span class="hud-label">高度</span>
        <span class="hud-value">{{ hud.alt }}<small>m</small></span>
      </div>
      <div class="hud-item">
        <span class="hud-label">航速</span>
        <span class="hud-value">{{ hud.spd }}<small>km/h</small></span>
      </div>
      <div class="hud-item">
        <span class="hud-label">航向</span>
        <span class="hud-value">{{ hud.dir }}</span>
      </div>
      <div class="hud-item">
        <span class="hud-label">模式</span>
        <span class="hud-value hud-mode">{{ autoPilot ? 'AUTO' : 'MANUAL' }}</span>
      </div>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar" v-if="!loading">
      <button :class="['tool-btn', { on: autoPilot }]" @click="toggleAuto">
        {{ autoPilot ? '🧭 自动巡航中' : '🧭 自动巡航' }}
      </button>
      <button :class="['tool-btn', { on: viewMode === 'cockpit' }]" @click="switchView">
        {{ viewMode === 'chase' ? '🎥 追尾视角' : viewMode === 'cockpit' ? '🎥 驾驶舱' : '🎥 自由视角' }}
      </button>
      <button class="tool-btn" @click="resetFlight">🔄 重置</button>
      <button :class="['tool-btn', { on: heatLayer }]" @click="toggleHeat">
        {{ heatLayer ? '📊 热力层开' : '📊 热力层关' }}
      </button>
    </div>

    <!-- 热力层透明度滑块 -->
    <div class="heat-slider" v-if="!loading && heatLayer">
      <span class="slider-label">热力层</span>
      <input
        type="range"
        min="15"
        max="90"
        step="5"
        :value="heatOpacity"
        @input="applyHeatOpacity(($event.target as HTMLInputElement).value)"
      />
      <span class="slider-value">{{ heatOpacity }}%</span>
    </div>

    <!-- 高度图例（热力层开启时显示） -->
    <div class="alt-legend" v-if="!loading && heatLayer">
      <span class="legend-label legend-top">高</span>
      <div class="legend-bar"></div>
      <span class="legend-label legend-bot">低</span>
      <span class="legend-value legend-max">{{ altRange.max }}m</span>
      <span class="legend-value legend-min">{{ altRange.min }}m</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from "vue";
import * as THREE from "three";
import { OrbitControls } from "three/examples/jsm/controls/OrbitControls.js";

const containerRef = ref<HTMLDivElement | null>(null);
const loading = ref(true);
const autoPilot = ref(true);
const viewMode = ref<"chase" | "cockpit" | "free">("chase");
const heatLayer = ref(true);
const heatOpacity = ref(55);
let heatOverlay: THREE.Mesh | null = null;
const altRange = reactive({ min: 0, max: 0 });

const hud = reactive({ alt: "0", spd: "0", dir: "N" });

// ========== 场景引用 ==========
let scene: THREE.Scene;
let camera: THREE.PerspectiveCamera;
let renderer: THREE.WebGLRenderer;
let controls: OrbitControls;
let animationId = 0;
const clock = new THREE.Clock();
let roObserver: ResizeObserver | null = null;
let disposed = false;

// ========== 地形参数 ==========
const TERRAIN_SIZE = 480;
const TERRAIN_SEG = 220;
const SEA_LEVEL = 0;

// ========== FBM 噪声（哈希值噪声 + 多倍频） ==========
function hash2(ix: number, iy: number): number {
  let n = ix * 374761393 + iy * 668265263;
  n = (n ^ (n >> 13)) * 1274126177;
  return ((n ^ (n >> 16)) >>> 0) / 4294967295;
}
function smooth(t: number): number {
  return t * t * (3 - 2 * t);
}
function valueNoise(x: number, y: number): number {
  const ix = Math.floor(x), iy = Math.floor(y);
  const fx = smooth(x - ix), fy = smooth(y - iy);
  const a = hash2(ix, iy), b = hash2(ix + 1, iy);
  const c = hash2(ix, iy + 1), d = hash2(ix + 1, iy + 1);
  return a + (b - a) * fx + (c - a) * fy + (a - b - c + d) * fx * fy;
}
function fbm(x: number, y: number, octaves = 5): number {
  let sum = 0, amp = 0.5, freq = 1, norm = 0;
  for (let i = 0; i < octaves; i++) {
    sum += valueNoise(x * freq, y * freq) * amp;
    norm += amp;
    amp *= 0.5;
    freq *= 2.07;
  }
  return sum / norm;
}

// 地形高度采样（与网格生成一致，供碰撞/巡航用）
function sampleHeight(x: number, z: number): number {
  let nx = x / TERRAIN_SIZE + 0.5;
  let nz = z / TERRAIN_SIZE + 0.5;
  // 域扭曲（domain warp）：让山体走向产生自然"侵蚀流线"，不再是均匀云絮状
  const wx = fbm(nx * 1.3 + 31.7, nz * 1.3 + 11.3, 3) - 0.5;
  const wz = fbm(nx * 1.3 + 71.3, nz * 1.3 + 53.9, 3) - 0.5;
  nx += wx * 0.42; nz += wz * 0.42;
  // 脊线噪声：山脉感
  const ridge = 1 - Math.abs(fbm(nx * 3.2, nz * 3.2, 4) * 2 - 1);
  const base = fbm(nx * 2.4 + 100, nz * 2.4 + 100, 5);
  let h = Math.pow(base, 1.6) * 34 + Math.pow(ridge, 2.2) * 30 - 8;
  // 边缘压低成海
  const edge = Math.min(nx, 1 - nx, nz, 1 - nz) * 6;
  h *= THREE.MathUtils.clamp(edge, 0, 1);
  return h;
}

// ========== 海拔分层着色（连续 colormap · matplotlib terrain 风格） ==========
// 深海蓝 → 海蓝 → 青 → 草绿 → 黄绿 → 金黄 → 橙 → 山顶红
const RAINBOW_CMAP: THREE.Color[] = [
  new THREE.Color("#03286b"), // 深海蓝
  new THREE.Color("#0a4cc6"), // 海蓝
  new THREE.Color("#07b9d6"), // 青（浅水）
  new THREE.Color("#1bc46a"), // 草绿
  new THREE.Color("#9ee02b"), // 黄绿
  new THREE.Color("#ffd400"), // 金黄
  new THREE.Color("#ff7a1a"), // 橙
  new THREE.Color("#ff2d2d"), // 山顶红
];

// ========== 写实地形色板（低饱和自然色） ==========
const REAL_CMAP = {
  DEEP:    new THREE.Color("#16323e"), // 深水
  SHALLOW: new THREE.Color("#2d5a6b"), // 浅水
  WET:     new THREE.Color("#9d8f74"), // 潮湿泥滩
  SAND:    new THREE.Color("#c2b28d"), // 沙滩
  GRASS:   new THREE.Color("#67744a"), // 草地（哑绿）
  SHRUB:   new THREE.Color("#4e5a3a"), // 灌木
  FOREST:  new THREE.Color("#39472e"), // 深林
  ROCK:    new THREE.Color("#6b6258"), // 岩石灰褐
  ROCK2:   new THREE.Color("#7b7268"), // 高山岩
  SNOW:    new THREE.Color("#edeff0"), // 雪
};
const clamp01 = (v: number) => THREE.MathUtils.clamp(v, 0, 1);

// 写实着色：海拔分带 + 坡度混岩（陡坡露岩/缓坡长植被）+ 噪声扰动雪线 + 细节斑驳
function colorRealistic(h: number, slope: number, jitter: number): THREE.Color {
  const c = new THREE.Color();
  if (h < SEA_LEVEL) {
    c.copy(REAL_CMAP.DEEP).lerp(REAL_CMAP.SHALLOW, clamp01((h + 8) / 8));
  } else {
    // 基础海拔带（缓变过渡）
    if (h < 0.9) c.copy(REAL_CMAP.WET).lerp(REAL_CMAP.SAND, h / 0.9);
    else if (h < 6)   c.copy(REAL_CMAP.SAND).lerp(REAL_CMAP.GRASS, (h - 0.9) / 5.1);
    else if (h < 14)  c.copy(REAL_CMAP.GRASS).lerp(REAL_CMAP.SHRUB, (h - 6) / 8);
    else if (h < 22)  c.copy(REAL_CMAP.SHRUB).lerp(REAL_CMAP.FOREST, (h - 14) / 8);
    else if (h < 30)  c.copy(REAL_CMAP.FOREST).lerp(REAL_CMAP.ROCK, (h - 22) / 8);
    else {
      // 雪线由噪声扰动，不再是完美等高线
      const snowLine = 34 + (jitter - 0.5) * 7;
      if (h < snowLine) c.copy(REAL_CMAP.ROCK).lerp(REAL_CMAP.ROCK2, clamp01((h - 30) / Math.max(snowLine - 30, 1)));
      else c.copy(REAL_CMAP.ROCK2).lerp(REAL_CMAP.SNOW, clamp01((h - snowLine) / 3));
    }
    // 坡度因子：陡坡植被无法附着 → 露出岩石
    if (h > 1.2) {
      const rockMix = clamp01((slope - 0.6) * 1.7);
      if (rockMix > 0) c.lerp(REAL_CMAP.ROCK, rockMix * (h > 26 ? 0.45 : 0.9));
    }
  }
  // 高频细节斑驳（模拟草皮/碎石混杂）
  c.offsetHSL((jitter - 0.5) * 0.02, (jitter - 0.5) * 0.05, (jitter - 0.5) * 0.09);
  return c;
}

// 连续彩虹配色：在 [minH, maxH] 上做归一化插值
function colorRainbow(h: number, jitter: number, minH: number, maxH: number): THREE.Color {
  const span = Math.max(maxH - minH, 0.001);
  const t = THREE.MathUtils.clamp((h - minH) / span, 0, 1);
  const k = t * (RAINBOW_CMAP.length - 1);
  const i = Math.floor(k);
  const f = k - i;
  const c = new THREE.Color();
  if (i >= RAINBOW_CMAP.length - 1) c.copy(RAINBOW_CMAP[RAINBOW_CMAP.length - 1]);
  else c.copy(RAINBOW_CMAP[i]).lerp(RAINBOW_CMAP[i + 1], f);
  // 轻微 hue/lightness jitter 让色彩"活"起来，但不过分
  c.offsetHSL((jitter - 0.5) * 0.018, 0, (jitter - 0.5) * 0.06);
  return c;
}

// 底座永远用写实配色；彩虹高程数据由独立热力层叠加（数字孪生风格）
let altitudeRange: { min: number; max: number } = { min: -8, max: 40 };
function colorByHeight(h: number, slope: number, jitter: number): THREE.Color {
  return colorRealistic(h, slope, jitter);
}

// ========== 飞行状态 ==========
const plane = new THREE.Group();
const flight = {
  pos: new THREE.Vector3(0, 40, -150),
  yaw: 0,          // 航向角
  pitch: 0,        // 俯仰角
  roll: 0,         // 滚转（视觉倾斜）
  speed: 42,       // m/s
  cruiseAlt: 46,   // 巡航目标高度（相对地形）
};
let propeller: THREE.Object3D | null = null;
const keys: Record<string, boolean> = {};

// 巡航航点
let waypoint: THREE.Vector3 | null = null;
let waypointMarker: THREE.Group | null = null;
const circleCenter = new THREE.Vector2(0, 0);

// 云
const clouds: THREE.Group[] = [];

// 临时变量
const _fwd = new THREE.Vector3();
const _target = new THREE.Vector3();
const _camPos = new THREE.Vector3();

// ========== 构建地形 ==========
function buildTerrain(): THREE.Mesh {
  const geo = new THREE.PlaneGeometry(TERRAIN_SIZE, TERRAIN_SIZE, TERRAIN_SEG, TERRAIN_SEG);
  geo.rotateX(-Math.PI / 2);
  const pos = geo.attributes.position as THREE.BufferAttribute;

  // 第一遍：采样所有顶点高度，算 min/max
  let minH = Infinity, maxH = -Infinity;
  for (let i = 0; i < pos.count; i++) {
    const x = pos.getX(i), z = pos.getZ(i);
    const h = sampleHeight(x, z);
    pos.setY(i, h);
    if (h < minH) minH = h;
    if (h > maxH) maxH = h;
  }
  // 限制色彩映射到「主要陆地」区间（去掉边缘深海异常值），让陆地色阶更鲜明
  // 取 5%~98% 分位做软裁剪
  const lo = minH + (maxH - minH) * 0.05;
  const hi = minH + (maxH - minH) * 0.98;
  altitudeRange = { min: lo, max: hi };
  altRange.min = Math.round(lo);
  altRange.max = Math.round(hi);

  // 第二遍：上色（含坡度计算——直接用规则网格的邻点高度差，零额外噪声采样）
  paintTerrain(geo);
  geo.computeVertexNormals();
  const mat = new THREE.MeshStandardMaterial({
    vertexColors: true,
    roughness: 0.95,
    metalness: 0.02,
    flatShading: true,
  });
  const mesh = new THREE.Mesh(geo, mat);
  mesh.receiveShadow = true;
  mesh.name = "terrain";
  return mesh;
}

// 依据当前 palette 对 PlaneGeometry 顶点上色（buildTerrain / togglePalette 共用）
function paintTerrain(geo: THREE.PlaneGeometry) {
  const pos = geo.attributes.position as THREE.BufferAttribute;
  const rowLen = TERRAIN_SEG + 1;
  const gridStep = TERRAIN_SIZE / TERRAIN_SEG;
  let colors = geo.attributes.color
    ? (geo.attributes.color.array as Float32Array)
    : new Float32Array(pos.count * 3);
  for (let i = 0; i < pos.count; i++) {
    const x = pos.getX(i), z = pos.getZ(i);
    const h = pos.getY(i);
    // 邻点索引（规则网格行主序；边界用自身）
    const iRight = i % rowLen === TERRAIN_SEG ? i : i + 1;
    const iDown = i + rowLen < pos.count ? i + rowLen : i;
    const dhdx = (pos.getY(iRight) - h) / gridStep;
    const dhdz = (pos.getY(iDown) - h) / gridStep;
    const slope = Math.sqrt(dhdx * dhdx + dhdz * dhdz);
    const c = colorByHeight(h, slope, hash2(Math.round(x * 3), Math.round(z * 3)));
    colors[i * 3] = c.r; colors[i * 3 + 1] = c.g; colors[i * 3 + 2] = c.b;
  }
  if (!geo.attributes.color) geo.setAttribute("color", new THREE.BufferAttribute(colors, 3));
  else (geo.attributes.color as THREE.BufferAttribute).needsUpdate = true;
}

// ========== 数据热力层（数字孪生：写实底座 + 彩虹高程叠加） ==========
function buildHeatOverlay(terrainGeo: THREE.PlaneGeometry): THREE.Mesh {
  const geo = terrainGeo.clone();
  const pos = geo.attributes.position as THREE.BufferAttribute;
  const colors = new Float32Array(pos.count * 3);
  for (let i = 0; i < pos.count; i++) {
    const x = pos.getX(i), z = pos.getZ(i);
    const h = pos.getY(i);
    const c = colorRainbow(h, hash2(Math.round(x * 3), Math.round(z * 3)), altitudeRange.min, altitudeRange.max);
    colors[i * 3] = c.r; colors[i * 3 + 1] = c.g; colors[i * 3 + 2] = c.b;
  }
  geo.setAttribute("color", new THREE.BufferAttribute(colors, 3));
  // MeshBasicMaterial 不受光照影响 → 呈现"数据层"的自发光观感
  const mat = new THREE.MeshBasicMaterial({
    vertexColors: true,
    transparent: true,
    opacity: heatOpacity.value / 100,
    depthWrite: false,
    polygonOffset: true,
    polygonOffsetFactor: -2,
    polygonOffsetUnits: -2,
  });
  const mesh = new THREE.Mesh(geo, mat);
  mesh.position.y = 0.3;
  mesh.renderOrder = 2;
  mesh.name = "heatOverlay";
  mesh.visible = heatLayer.value;
  return mesh;
}

function buildWater(): THREE.Mesh {
  const geo = new THREE.PlaneGeometry(TERRAIN_SIZE * 1.05, TERRAIN_SIZE * 1.05);
  geo.rotateX(-Math.PI / 2);
  const mat = new THREE.MeshStandardMaterial({
    color: 0x2d5a6b,
    transparent: true,
    opacity: 0.78,
    roughness: 0.12,
    metalness: 0.7,
  });
  const mesh = new THREE.Mesh(geo, mat);
  mesh.position.y = SEA_LEVEL - 0.35;
  mesh.name = "water";
  return mesh;
}

// ========== 程序化低多边形飞机 ==========
function buildPlane(): THREE.Group {
  const g = new THREE.Group();

  const bodyMat = new THREE.MeshStandardMaterial({ color: 0xf0f4f8, roughness: 0.4, metalness: 0.3, flatShading: true });
  const accentMat = new THREE.MeshStandardMaterial({ color: 0xe74c3c, roughness: 0.4, metalness: 0.3, flatShading: true });
  const darkMat = new THREE.MeshStandardMaterial({ color: 0x2c3e50, roughness: 0.5, metalness: 0.4, flatShading: true });

  // 机身（锥头 + 圆柱 + 锥尾）
  const nose = new THREE.Mesh(new THREE.ConeGeometry(0.7, 1.6, 10), accentMat);
  nose.rotation.x = -Math.PI / 2;
  nose.position.z = 2.6;
  const fuselage = new THREE.Mesh(new THREE.CylinderGeometry(0.7, 0.5, 3.4, 10), bodyMat);
  fuselage.rotation.x = Math.PI / 2;
  fuselage.position.z = 0.3;
  const tailCone = new THREE.Mesh(new THREE.ConeGeometry(0.5, 1.6, 10), bodyMat);
  tailCone.rotation.x = Math.PI / 2;
  tailCone.position.z = -2.2;

  // 主翼
  const wing = new THREE.Mesh(new THREE.BoxGeometry(7.2, 0.16, 1.5), bodyMat);
  wing.position.set(0, 0.15, 0.35);
  const wingTipL = new THREE.Mesh(new THREE.BoxGeometry(0.5, 0.18, 1.1), accentMat);
  wingTipL.position.set(-3.6, 0.15, 0.35);
  const wingTipR = wingTipL.clone();
  wingTipR.position.x = 3.6;

  // 尾翼
  const hStab = new THREE.Mesh(new THREE.BoxGeometry(2.6, 0.12, 0.8), bodyMat);
  hStab.position.set(0, 0.25, -2.4);
  const vStab = new THREE.Mesh(new THREE.BoxGeometry(0.12, 1.2, 0.9), accentMat);
  vStab.position.set(0, 0.85, -2.4);

  // 螺旋桨
  const propGroup = new THREE.Group();
  const spinner = new THREE.Mesh(new THREE.ConeGeometry(0.28, 0.7, 8), darkMat);
  spinner.rotation.x = -Math.PI / 2;
  const blade1 = new THREE.Mesh(new THREE.BoxGeometry(0.14, 2.8, 0.06), darkMat);
  const blade2 = blade1.clone();
  blade2.rotation.z = Math.PI / 2;
  propGroup.add(spinner, blade1, blade2);
  propGroup.position.z = 3.5;
  propeller = propGroup;

  // 座舱
  const canopy = new THREE.Mesh(
    new THREE.SphereGeometry(0.45, 10, 8),
    new THREE.MeshStandardMaterial({ color: 0x4da6ff, roughness: 0.1, metalness: 0.6, transparent: true, opacity: 0.85 })
  );
  canopy.scale.set(1, 0.75, 1.6);
  canopy.position.set(0, 0.55, 1.1);

  g.add(nose, fuselage, tailCone, wing, wingTipL, wingTipR, hStab, vStab, propGroup, canopy);
  g.traverse((o) => { o.castShadow = true; });
  return g;
}

// ========== 云朵 ==========
function buildClouds() {
  const cloudMat = new THREE.MeshStandardMaterial({ color: 0xffffff, roughness: 1, transparent: true, opacity: 0.88, flatShading: true });
  for (let i = 0; i < 14; i++) {
    const cloud = new THREE.Group();
    const puffs = 3 + Math.floor(Math.random() * 3);
    for (let p = 0; p < puffs; p++) {
      const r = 3 + Math.random() * 4;
      const puff = new THREE.Mesh(new THREE.IcosahedronGeometry(r, 0), cloudMat);
      puff.position.set((Math.random() - 0.5) * 12, (Math.random() - 0.5) * 2.5, (Math.random() - 0.5) * 8);
      puff.scale.y = 0.55;
      cloud.add(puff);
    }
    const angle = Math.random() * Math.PI * 2;
    const dist = 60 + Math.random() * 160;
    cloud.position.set(Math.cos(angle) * dist, 55 + Math.random() * 45, Math.sin(angle) * dist);
    cloud.userData.speed = 1.2 + Math.random() * 1.6;
    scene.add(cloud);
    clouds.push(cloud);
  }
}

// ========== 航点标记 ==========
function setWaypoint(x: number, z: number) {
  clearWaypoint();
  const h = sampleHeight(x, z);
  const marker = new THREE.Group();

  const ringGeo = new THREE.RingGeometry(2.2, 3.2, 40);
  const ringMat = new THREE.MeshBasicMaterial({ color: 0x4da6ff, transparent: true, opacity: 0.8, side: THREE.DoubleSide, depthWrite: false });
  const ring = new THREE.Mesh(ringGeo, ringMat);
  ring.rotation.x = -Math.PI / 2;
  ring.position.y = 0.15;
  marker.add(ring);

  const beamGeo = new THREE.CylinderGeometry(0.35, 0.35, 60, 8, 1, true);
  const beamMat = new THREE.MeshBasicMaterial({ color: 0x4da6ff, transparent: true, opacity: 0.3, side: THREE.DoubleSide, depthWrite: false });
  const beam = new THREE.Mesh(beamGeo, beamMat);
  beam.position.y = 30;
  marker.add(beam);

  marker.position.set(x, Math.max(h, SEA_LEVEL), z);
  scene.add(marker);
  waypointMarker = marker;
  waypoint = new THREE.Vector3(x, 0, z);
  autoPilot.value = true;
}

function clearWaypoint() {
  if (waypointMarker) {
    scene.remove(waypointMarker);
    waypointMarker.traverse((o) => {
      const m = o as THREE.Mesh;
      if (m.isMesh) { m.geometry.dispose(); (m.material as THREE.Material).dispose(); }
    });
    waypointMarker = null;
  }
  waypoint = null;
}

// ========== 交互 ==========
const raycaster = new THREE.Raycaster();
const ndc = new THREE.Vector2();

function onPointerDown(e: PointerEvent) {
  if (!renderer || !containerRef.value) return;
  const rect = renderer.domElement.getBoundingClientRect();
  if (e.clientX < rect.left || e.clientX > rect.right || e.clientY < rect.top || e.clientY > rect.bottom) return;
  ndc.x = ((e.clientX - rect.left) / rect.width) * 2 - 1;
  ndc.y = -((e.clientY - rect.top) / rect.height) * 2 + 1;
  raycaster.setFromCamera(ndc, camera);
  const terrain = scene.getObjectByName("terrain");
  if (!terrain) return;
  const hits = raycaster.intersectObject(terrain, false);
  if (hits.length > 0 && viewMode.value !== "free") {
    setWaypoint(hits[0].point.x, hits[0].point.z);
  }
}

function onKeyDown(e: KeyboardEvent) {
  keys[e.code] = true;
  if (["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", "Space"].includes(e.code)) e.preventDefault();
}
function onKeyUp(e: KeyboardEvent) {
  keys[e.code] = false;
}

// ========== 视角 ==========
function switchView() {
  viewMode.value = viewMode.value === "chase" ? "cockpit" : viewMode.value === "cockpit" ? "free" : "chase";
  if (viewMode.value === "free") {
    controls.enabled = true;
    controls.target.copy(flight.pos);
    camera.position.copy(flight.pos).add(new THREE.Vector3(18, 10, 18));
  } else {
    controls.enabled = false;
  }
}

function toggleAuto() {
  autoPilot.value = !autoPilot.value;
  if (autoPilot.value) {
    flight.pitch *= 0.3;
    flight.roll *= 0.3;
  } else {
    clearWaypoint();
  }
}

function resetFlight() {
  flight.pos.set(0, 40, -150);
  flight.yaw = 0; flight.pitch = 0; flight.roll = 0;
  flight.speed = 42;
  clearWaypoint();
  viewMode.value = "chase";
  if (controls) controls.enabled = false;
  autoPilot.value = true;
}

function toggleHeat() {
  heatLayer.value = !heatLayer.value;
  if (heatOverlay) heatOverlay.visible = heatLayer.value;
}

function applyHeatOpacity(v: string | number) {
  heatOpacity.value = Number(v);
  if (heatOverlay) (heatOverlay.material as THREE.MeshBasicMaterial).opacity = heatOpacity.value / 100;
}

// ========== 飞行逻辑 ==========
function tickFlight(dt: number, elapsed: number) {
  const groundH = sampleHeight(flight.pos.x, flight.pos.z);
  const terrainFollow = Math.max(groundH, SEA_LEVEL) + flight.cruiseAlt;

  if (autoPilot.value) {
    // ===== 自动巡航 =====
    let targetX: number, targetZ: number;
    if (waypoint) {
      targetX = waypoint.x; targetZ = waypoint.z;
      // 到达航点后回到环绕巡航
      const dwp = Math.hypot(targetX - flight.pos.x, targetZ - flight.pos.z);
      if (dwp < 15) clearWaypoint();
    } else {
      // 大圆环绕
      const ang = elapsed * 0.052;
      targetX = circleCenter.x + Math.cos(ang) * 150;
      targetZ = circleCenter.y + Math.sin(ang) * 150;
    }
    const wantYaw = Math.atan2(targetX - flight.pos.x, targetZ - flight.pos.z);
    let dyaw = wantYaw - flight.yaw;
    while (dyaw > Math.PI) dyaw -= Math.PI * 2;
    while (dyaw < -Math.PI) dyaw += Math.PI * 2;
    const turnRate = THREE.MathUtils.clamp(dyaw * 1.6, -0.55, 0.55);
    flight.yaw += turnRate * dt;
    // 转弯压坡度
    flight.roll = THREE.MathUtils.lerp(flight.roll, THREE.MathUtils.clamp(turnRate * 1.3, -0.7, 0.7), dt * 2.5);
    // 高度趋近巡航线，前方地形前瞻
    const aheadH = sampleHeight(
      flight.pos.x + Math.sin(flight.yaw) * 40,
      flight.pos.z + Math.cos(flight.yaw) * 40
    );
    const wantAlt = Math.max(terrainFollow, Math.max(aheadH, SEA_LEVEL) + flight.cruiseAlt * 0.9);
    flight.pos.y = THREE.MathUtils.lerp(flight.pos.y, wantAlt, dt * 0.8);
    flight.pitch = THREE.MathUtils.lerp(flight.pitch, THREE.MathUtils.clamp((wantAlt - flight.pos.y) * 0.02, -0.35, 0.35), dt * 2);
    flight.speed = THREE.MathUtils.lerp(flight.speed, 46, dt);
  } else {
    // ===== 手动驾驶 =====
    const turn = (keys["KeyA"] ? 1 : 0) - (keys["KeyD"] ? 1 : 0);
    flight.yaw += turn * 1.1 * dt;
    flight.roll = THREE.MathUtils.lerp(flight.roll, turn * 0.65, dt * 4);

    const pitchIn = (keys["ArrowUp"] ? 1 : 0) - (keys["ArrowDown"] ? 1 : 0);
    flight.pitch = THREE.MathUtils.clamp(flight.pitch + pitchIn * 1.2 * dt, -0.6, 0.6);
    if (!pitchIn) flight.pitch = THREE.MathUtils.lerp(flight.pitch, 0, dt * 1.5);

    const throttle = (keys["KeyW"] ? 1 : 0) - (keys["KeyS"] ? 1 : 0);
    flight.speed = THREE.MathUtils.clamp(flight.speed + throttle * 22 * dt, 16, 95);

    _fwd.set(Math.sin(flight.yaw) * Math.cos(flight.pitch), Math.sin(flight.pitch), Math.cos(flight.yaw) * Math.cos(flight.pitch));
    flight.pos.y += _fwd.y * flight.speed * dt;

    // 地形防撞：最低高度钳制
    const minY = Math.max(groundH, SEA_LEVEL) + 4;
    if (flight.pos.y < minY) {
      flight.pos.y = minY;
      if (flight.pitch < 0) flight.pitch = 0;
    }
    if (flight.pos.y > 160) flight.pos.y = 160;
  }

  // 位置推进
  _fwd.set(Math.sin(flight.yaw), 0, Math.cos(flight.yaw));
  flight.pos.x += _fwd.x * flight.speed * dt;
  flight.pos.z += _fwd.z * flight.speed * dt;

  // 边界回弹
  const half = TERRAIN_SIZE / 2 - 30;
  if (Math.abs(flight.pos.x) > half) { flight.pos.x = THREE.MathUtils.clamp(flight.pos.x, -half, half); flight.yaw += Math.PI * 0.5; if (waypoint) clearWaypoint(); }
  if (Math.abs(flight.pos.z) > half) { flight.pos.z = THREE.MathUtils.clamp(flight.pos.z, -half, half); flight.yaw += Math.PI * 0.5; if (waypoint) clearWaypoint(); }

  // 应用到模型（模型头朝 +Z）
  plane.position.copy(flight.pos);
  plane.rotation.set(0, 0, 0);
  plane.rotateY(flight.yaw);
  plane.rotateX(-flight.pitch);
  plane.rotateZ(-flight.roll);

  // 螺旋桨
  if (propeller) propeller.rotation.z += dt * (18 + flight.speed * 0.35);

  // ===== 相机 =====
  if (viewMode.value === "chase") {
    const back = 16, up = 6;
    _camPos.set(
      flight.pos.x - Math.sin(flight.yaw) * back,
      flight.pos.y + up,
      flight.pos.z - Math.cos(flight.yaw) * back
    );
    const groundClamp = Math.max(sampleHeight(_camPos.x, _camPos.z), SEA_LEVEL) + 3;
    if (_camPos.y < groundClamp) _camPos.y = groundClamp;
    camera.position.lerp(_camPos, 1 - Math.pow(0.001, dt));
    _target.copy(flight.pos).addScaledVector(_fwd, 12);
    camera.lookAt(_target);
  } else if (viewMode.value === "cockpit") {
    camera.position.copy(flight.pos).add(new THREE.Vector3(0, 0.9, 0.5).applyEuler(plane.rotation));
    _target.copy(flight.pos).addScaledVector(_fwd, 30);
    _target.y += flight.pitch * -10;
    camera.lookAt(_target);
  } else {
    controls.target.lerp(flight.pos, dt * 3);
    controls.update();
  }

  // ===== HUD =====
  hud.alt = Math.max(0, Math.round(flight.pos.y - Math.max(groundH, SEA_LEVEL))).toString();
  hud.spd = Math.round(flight.speed * 3.6).toString();
  const deg = ((THREE.MathUtils.radToDeg(flight.yaw) % 360) + 360) % 360;
  const dirs = ["N", "NE", "E", "SE", "S", "SW", "W", "NW"];
  hud.dir = dirs[Math.round(deg / 45) % 8] + " " + Math.round(deg) + "°";
}

// ========== 初始化 ==========
function init() {
  if (!containerRef.value) return;
  const w = containerRef.value.clientWidth;
  const h = containerRef.value.clientHeight;

  scene = new THREE.Scene();
  scene.background = new THREE.Color(0x87b8e8);
  scene.fog = new THREE.Fog(0xa8c8e8, 180, 620);

  camera = new THREE.PerspectiveCamera(60, w / h, 0.1, 1500);
  camera.position.set(0, 50, -170);

  renderer = new THREE.WebGLRenderer({ antialias: true });
  renderer.setSize(w, h);
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
  renderer.shadowMap.enabled = true;
  renderer.shadowMap.type = THREE.PCFSoftShadowMap;
  containerRef.value.appendChild(renderer.domElement);

  // 灯光
  const sun = new THREE.DirectionalLight(0xfff2d9, 1.6);
  sun.position.set(120, 180, 80);
  sun.castShadow = true;
  sun.shadow.mapSize.set(2048, 2048);
  sun.shadow.camera.left = -260; sun.shadow.camera.right = 260;
  sun.shadow.camera.top = 260; sun.shadow.camera.bottom = -260;
  sun.shadow.camera.far = 600;
  scene.add(sun);
  scene.add(new THREE.HemisphereLight(0xbcd8f5, 0x6b7f6b, 0.75));

  // 地形（写实底座）+ 数据热力层 + 水面 + 云
  const terrainMesh = buildTerrain();
  scene.add(terrainMesh);
  heatOverlay = buildHeatOverlay(terrainMesh.geometry as THREE.PlaneGeometry);
  scene.add(heatOverlay);
  scene.add(buildWater());
  buildClouds();

  // 飞机（跨挂载复用同一 Group，先清空旧模型）
  plane.clear();
  const planeModel = buildPlane();
  plane.add(planeModel);
  scene.add(plane);
  // 重置飞行状态
  flight.pos.set(0, 40, -150);
  flight.yaw = 0; flight.pitch = 0; flight.roll = 0;
  flight.speed = 42;

  // OrbitControls（默认关闭，自由视角时启用）
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true;
  controls.dampingFactor = 0.06;
  controls.minDistance = 8;
  controls.maxDistance = 400;
  controls.enabled = false;

  // 事件
  renderer.domElement.addEventListener("pointerdown", onPointerDown);
  window.addEventListener("keydown", onKeyDown);
  window.addEventListener("keyup", onKeyUp);

  roObserver = new ResizeObserver(() => {
    if (!containerRef.value || disposed) return;
    const cw = containerRef.value.clientWidth;
    const ch = containerRef.value.clientHeight;
    if (cw === 0 || ch === 0) return;
    camera.aspect = cw / ch;
    camera.updateProjectionMatrix();
    renderer.setSize(cw, ch);
  });
  roObserver.observe(containerRef.value);

  clock.getDelta(); // 丢弃首帧
  loading.value = false;

  function animate() {
    if (disposed) return;
    animationId = requestAnimationFrame(animate);
    const dt = Math.min(clock.getDelta(), 0.05);
    const elapsed = clock.getElapsedTime();

    tickFlight(dt, elapsed);

    // 云漂移
    for (const c of clouds) {
      c.position.x += c.userData.speed * dt;
      if (c.position.x > TERRAIN_SIZE / 2 + 40) c.position.x = -TERRAIN_SIZE / 2 - 40;
    }

    // 航点呼吸动画
    if (waypointMarker) {
      const s = 1 + Math.sin(elapsed * 4) * 0.15;
      waypointMarker.children[0].scale.set(s, s, 1);
    }

    renderer.render(scene, camera);
  }
  animate();
}

function dispose() {
  disposed = true;
  cancelAnimationFrame(animationId);
  roObserver?.disconnect();
  window.removeEventListener("keydown", onKeyDown);
  window.removeEventListener("keyup", onKeyUp);
  if (renderer?.domElement) renderer.domElement.removeEventListener("pointerdown", onPointerDown);
  controls?.dispose();
  clearWaypoint();
  for (const c of clouds) scene?.remove(c);
  clouds.length = 0;
  heatOverlay = null;
  scene?.traverse((o) => {
    const m = o as THREE.Mesh;
    if (m.isMesh) {
      m.geometry?.dispose();
      const mat = m.material;
      if (Array.isArray(mat)) mat.forEach((x) => x.dispose());
      else mat?.dispose();
    }
  });
  renderer?.dispose();
  if (renderer?.domElement?.parentElement) renderer.domElement.parentElement.removeChild(renderer.domElement);
}

onMounted(() => init());
onUnmounted(() => dispose());
</script>

<style scoped>
.terrain-demo {
  width: 100%;
  height: 100%;
  position: relative;
  background: #87b8e8;
  cursor: crosshair;
}

.terrain-loading {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  background: rgba(10, 16, 28, 0.9);
  z-index: 10;
  color: #cfe4ff;
  font-size: 0.9rem;
}

.loading-spinner {
  width: 36px;
  height: 36px;
  border: 3px solid rgba(77, 166, 255, 0.2);
  border-top-color: #4da6ff;
  border-radius: 50%;
  animation: spin 0.9s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.hud {
  position: absolute;
  top: 12px;
  left: 12px;
  display: flex;
  gap: 8px;
  z-index: 5;
}

.hud-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 64px;
  padding: 6px 10px;
  background: rgba(8, 16, 30, 0.55);
  backdrop-filter: blur(6px);
  border: 1px solid rgba(77, 166, 255, 0.35);
  border-radius: 8px;
}

.hud-label {
  font-size: 0.64rem;
  color: #8fb8e0;
  letter-spacing: 0.08em;
}

.hud-value {
  font-size: 0.95rem;
  font-weight: 700;
  color: #eaf4ff;
  font-variant-numeric: tabular-nums;
}

.hud-value small {
  font-size: 0.6rem;
  font-weight: 400;
  margin-left: 2px;
  color: #8fb8e0;
}

.hud-mode {
  color: #4da6ff;
}

.toolbar {
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 8px;
  z-index: 5;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.tool-btn {
  padding: 6px 12px;
  font-size: 0.74rem;
  border: 1px solid rgba(77, 166, 255, 0.4);
  border-radius: 8px;
  background: rgba(8, 16, 30, 0.55);
  backdrop-filter: blur(6px);
  color: #cfe4ff;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.tool-btn:hover {
  border-color: #4da6ff;
  color: #fff;
}

.tool-btn.on {
  background: linear-gradient(135deg, rgba(77, 166, 255, 0.35), rgba(168, 85, 247, 0.3));
  border-color: #4da6ff;
  color: #fff;
}

/* ===== 热力层透明度滑块 ===== */
.heat-slider {
  position: absolute;
  right: 24px;
  bottom: 24px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  background: rgba(8, 16, 30, 0.55);
  backdrop-filter: blur(6px);
  border: 1px solid rgba(77, 166, 255, 0.4);
  border-radius: 10px;
  z-index: 6;
}

.heat-slider .slider-label {
  font-size: 0.72rem;
  color: #cfe4ff;
  white-space: nowrap;
}

.heat-slider .slider-value {
  font-size: 0.72rem;
  color: #fff;
  min-width: 34px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.heat-slider input[type="range"] {
  width: 110px;
  height: 4px;
  appearance: none;
  -webkit-appearance: none;
  border-radius: 2px;
  background: linear-gradient(90deg, #4da6ff, #a855f7);
  cursor: pointer;
}

.heat-slider input[type="range"]::-webkit-slider-thumb {
  appearance: none;
  -webkit-appearance: none;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #fff;
  border: 2px solid #4da6ff;
}

/* ===== 高度图例 ===== */
.alt-legend {
  position: absolute;
  right: 24px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  pointer-events: none;
  user-select: none;
  z-index: 5;
}

.alt-legend .legend-bar {
  width: 14px;
  height: 200px;
  border-radius: 7px;
  border: 1px solid rgba(255, 255, 255, 0.45);
  background: linear-gradient(
    to bottom,
    #ff2d2d 0%,
    #ff7a1a 15%,
    #ffd400 30%,
    #9ee02b 45%,
    #1bc46a 60%,
    #07b9d6 75%,
    #0a4cc6 90%,
    #03286b 100%
  );
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.35), inset 0 0 0 1px rgba(0, 0, 0, 0.2);
}

.alt-legend .legend-label {
  font-family: "Consolas", "Monaco", monospace;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 1px;
  color: #fff;
  text-shadow: 0 1px 4px rgba(0, 0, 0, 0.6);
}

.alt-legend .legend-value {
  font-family: "Consolas", "Monaco", monospace;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.92);
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.6);
  position: absolute;
  right: 22px;
}

.alt-legend .legend-max { top: 28px; }
.alt-legend .legend-min { top: calc(28px + 196px); }
</style>
