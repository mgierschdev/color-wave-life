const canvas = document.getElementById("life");
const ctx = canvas.getContext("2d");
const toggle = document.getElementById("toggle");
const reset = document.getElementById("reset");
const slower = document.getElementById("slower");
const faster = document.getElementById("faster");
const record = document.getElementById("record");
const patternSelect = document.getElementById("pattern");
const speedInput = document.getElementById("speed-input");
const statusNode = document.getElementById("status");
const patternLabelNode = document.getElementById("pattern-label");
const generationNode = document.getElementById("generation");
const speedNode = document.getElementById("speed");

const width = 240;
const height = 160;
const cellSize = 4;
const minSimulationFPS = 1;
const waveSpeed = 0.012;
const wavelength = 11;
const saturation = 85;
const brightness = 71;
const background = "#080a12";
const grid = "#141824";

let world = createEmptyWorld();
let next = createEmptyWorld();
let phase = 0;
let generation = 0;
let running = true;
let recorder = null;
let recordedChunks = [];
let lastFrameTime = 0;
let simulationAccumulator = 0;
let simulationFPS = 5;
let currentPattern = "mixed";

resetWorld();
updateSpeedLabel();
draw();
requestAnimationFrame(loop);

toggle.addEventListener("click", () => {
  running = !running;
  statusNode.textContent = running ? "running" : "paused";
  toggle.textContent = running ? "Pause" : "Resume";
});

reset.addEventListener("click", () => {
  resetWorld();
  draw();
});

patternSelect.addEventListener("change", () => {
  currentPattern = patternSelect.value;
  resetWorld();
  draw();
});

speedInput.addEventListener("change", () => {
  setSimulationFPS(speedInput.value);
});

slower.addEventListener("click", () => {
  setSimulationFPS(simulationFPS - 1);
});

faster.addEventListener("click", () => {
  setSimulationFPS(simulationFPS + 1);
});

record.addEventListener("click", () => {
  if (recorder && recorder.state === "recording") {
    return;
  }
  recordedChunks = [];
  recorder = new MediaRecorder(canvas.captureStream(20), { mimeType: "video/webm" });
  recorder.ondataavailable = (event) => {
    if (event.data.size > 0) {
      recordedChunks.push(event.data);
    }
  };
  recorder.onstop = () => {
    const blob = new Blob(recordedChunks, { type: "video/webm" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "color-wave-life.webm";
    a.click();
    setTimeout(() => URL.revokeObjectURL(url), 5000);
    record.textContent = "Record 12s";
  };
  recorder.start();
  record.textContent = "Recording...";
  setTimeout(() => recorder.stop(), 12000);
});

function loop(now) {
  if (!lastFrameTime) {
    lastFrameTime = now;
  }
  const dt = (now - lastFrameTime) / 1000;
  lastFrameTime = now;

  if (running) {
    simulationAccumulator += dt;
    const simulationStep = 1 / simulationFPS;
    while (simulationAccumulator >= simulationStep) {
      step();
      simulationAccumulator -= simulationStep;
    }
    phase += waveSpeed * dt * 60;
  }
  draw();
  requestAnimationFrame(loop);
}

function resetWorld() {
  world = createEmptyWorld();
  next = createEmptyWorld();
  const cx = Math.floor(width / 2);
  const cy = Math.floor(height / 2);
  if (currentPattern === "mixed") {
    seedMixedWorld(cx, cy);
    phase = 0;
    generation = 0;
    simulationAccumulator = 0;
    generationNode.textContent = `generation: ${generation}`;
    statusNode.textContent = running ? "running" : "paused";
    patternLabelNode.textContent = "pattern: mixed expansion";
    return;
  }
  const pattern = getPattern(currentPattern);
  const originX = cx + pattern.originX;
  const originY = cy + pattern.originY;
  for (const [dx, dy] of pattern.cells) {
    setAlive(world, originX + dx, originY + dy, true);
  }
  phase = 0;
  generation = 0;
  simulationAccumulator = 0;
  generationNode.textContent = `generation: ${generation}`;
  statusNode.textContent = running ? "running" : "paused";
  patternLabelNode.textContent = `pattern: ${pattern.label}`;
}

function seedMixedWorld(cx, cy) {
  const seeds = [
    { name: "spacefiller", x: Math.floor(width * 0.16), y: Math.floor(height * 0.18) },
    { name: "glidergun", x: Math.floor(width * 0.10), y: Math.floor(height * 0.62) },
    { name: "switchengine", x: Math.floor(width * 0.34), y: Math.floor(height * 0.22) },
    { name: "pulsar", x: Math.floor(width * 0.52), y: Math.floor(height * 0.20) },
    { name: "rpentomino", x: Math.floor(width * 0.74), y: Math.floor(height * 0.18) },
    { name: "acorn", x: Math.floor(width * 0.88), y: Math.floor(height * 0.30) },
    { name: "glider", x: Math.floor(width * 0.24), y: Math.floor(height * 0.46) },
    { name: "diehard", x: Math.floor(width * 0.46), y: Math.floor(height * 0.54) },
    { name: "lwss", x: Math.floor(width * 0.68), y: Math.floor(height * 0.58) },
    { name: "switchengine", x: Math.floor(width * 0.86), y: Math.floor(height * 0.62) },
    { name: "acorn", x: Math.floor(width * 0.22), y: Math.floor(height * 0.80) },
    { name: "rpentomino", x: Math.floor(width * 0.50), y: Math.floor(height * 0.80) },
    { name: "glider", x: Math.floor(width * 0.76), y: Math.floor(height * 0.80) },
  ];

  for (const seed of seeds) {
    const pattern = getPattern(seed.name);
    for (const [dx, dy] of pattern.cells) {
      setAlive(world, seed.x + dx, seed.y + dy, true);
    }
  }
}

function updateSpeedLabel() {
  speedNode.textContent = `speed: ${simulationFPS.toFixed(1)} gen/s`;
  speedInput.value = String(Math.round(simulationFPS));
}

function setSimulationFPS(value) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    speedInput.value = String(Math.round(simulationFPS));
    return;
  }
  simulationFPS = Math.max(minSimulationFPS, parsed);
  simulationAccumulator = 0;
  updateSpeedLabel();
}

function getPattern(name) {
  if (name === "acorn") {
    return {
      label: "acorn",
      originX: 0,
      originY: 0,
      cells: [
        [-3, 0],
        [-1, 1],
        [-3, 2], [-2, 2], [1, 2],
        [-3, 3],
      ],
    };
  }
  if (name === "spacefiller") {
    return {
      label: "spacefiller",
      originX: -24,
      originY: -13,
      cells: [
        [-4, -11], [-4, -10], [-4, -9], [0, -11], [0, -10], [0, -9],
        [-5, -10], [-1, -10],
        [-24, -9], [-23, -9], [-22, -9], [-21, -9],
        [-24, -8], [-21, -8],
        [-24, -3], [-21, -3],
        [-18, -3],
        [-17, -2], [-15, -2],
        [-24, -1], [-21, -1], [-19, -1], [-18, -1],
        [-17, -1], [-16, -1], [-15, -1], [-14, -1], [-13, -1],
        [-21, 0], [-20, 0],
        [-12, 0], [-11, 0],
        [-24, 1], [-23, 1], [-22, 1], [-21, 1],
        [-10, 1], [-9, 1], [-8, 1], [-7, 1],
        [-18, 2], [-12, 2], [-11, 2],
        [-5, 2], [-3, 2], [-2, 2], [0, 2],
        [-24, 3], [-12, 3], [-8, 3], [-5, 3], [-4, 3], [-3, 3], [-1, 3],
        [-24, 4], [-22, 4], [-20, 4], [-19, 4], [-18, 4], [-13, 4], [-7, 4],
        [-15, 5], [-11, 5], [-10, 5], [-9, 5], [-8, 5], [-5, 5], [-3, 5],
        [-18, 6], [-17, 6], [-15, 6], [-14, 6], [-10, 6], [-7, 6], [-5, 6], [-4, 6],
        [-3, 6], [-2, 6], [-1, 6], [0, 6],
        [-18, 7], [-16, 7], [-15, 7], [-14, 7], [-10, 7], [-9, 7], [-8, 7], [-5, 7],
        [-3, 7], [-1, 7],
        [-19, 8], [-18, 8], [-17, 8], [-15, 8], [-10, 8], [-5, 8], [-3, 8],
        [-2, 8], [-1, 8], [0, 8],
        [-21, 9], [-20, 9], [-19, 9], [-18, 9], [-14, 9], [-13, 9], [-12, 9],
        [-11, 9], [-10, 9],
        [-24, 10], [-23, 10], [-22, 10], [-21, 10],
        [-15, 11], [-14, 11], [-13, 11],
        [-18, 12],
        [-24, 13], [-23, 13], [-22, 13], [-21, 13],
      ],
    };
  }
  if (name === "rpentomino") {
    return {
      label: "r-pentomino",
      originX: 0,
      originY: 0,
      cells: [
        [0, -1], [1, -1],
        [-1, 0], [0, 0],
        [0, 1],
      ],
    };
  }
  if (name === "glider") {
    return {
      label: "glider",
      originX: 0,
      originY: 0,
      cells: [
        [0, -1],
        [1, 0],
        [-1, 1], [0, 1], [1, 1],
      ],
    };
  }
  if (name === "lwss") {
    return {
      label: "lwss",
      originX: 0,
      originY: 0,
      cells: [
        [-2, -1], [-1, -1], [0, -1], [1, -1],
        [-3, 0], [1, 0],
        [1, 1],
        [-3, 2], [0, 2],
      ],
    };
  }
  if (name === "diehard") {
    return {
      label: "diehard",
      originX: 0,
      originY: 0,
      cells: [
        [-3, 0],
        [-2, 0],
        [-2, 1],
        [2, 1],
        [3, -1], [3, 1],
        [4, 1],
      ],
    };
  }
  if (name === "pulsar") {
    return {
      label: "pulsar",
      originX: 0,
      originY: 0,
      cells: [
        [-4, -6], [-3, -6], [-2, -6], [2, -6], [3, -6], [4, -6],
        [-6, -4], [-1, -4], [1, -4], [6, -4],
        [-6, -3], [-1, -3], [1, -3], [6, -3],
        [-6, -2], [-1, -2], [1, -2], [6, -2],
        [-4, -1], [-3, -1], [-2, -1], [2, -1], [3, -1], [4, -1],
        [-4, 1], [-3, 1], [-2, 1], [2, 1], [3, 1], [4, 1],
        [-6, 2], [-1, 2], [1, 2], [6, 2],
        [-6, 3], [-1, 3], [1, 3], [6, 3],
        [-6, 4], [-1, 4], [1, 4], [6, 4],
        [-4, 6], [-3, 6], [-2, 6], [2, 6], [3, 6], [4, 6],
      ],
    };
  }
  if (name === "switchengine") {
    return {
      label: "switch engine",
      originX: 0,
      originY: 0,
      cells: [
        [-4, -2], [-3, -2],
        [-2, -1], [-4, -1],
        [-4, 0],
        [-2, 1], [0, 1], [1, 1],
        [2, 0],
        [0, -2],
      ],
    };
  }
  return {
    label: "glider gun",
    originX: -18,
    originY: -4,
    cells: [
      [0, 4], [0, 5], [1, 4], [1, 5],
      [10, 4], [10, 5], [10, 6],
      [11, 3], [11, 7],
      [12, 2], [12, 8],
      [13, 2], [13, 8],
      [14, 5],
      [15, 3], [15, 7],
      [16, 4], [16, 5], [16, 6],
      [17, 5],
      [20, 2], [20, 3], [20, 4],
      [21, 2], [21, 3], [21, 4],
      [22, 1], [22, 5],
      [24, 0], [24, 1], [24, 5], [24, 6],
      [34, 2], [34, 3], [35, 2], [35, 3],
    ],
  };
}

function step() {
  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const neighbors = neighborCount(world, x, y);
      const alive = getAlive(world, x, y);
      next[index(x, y)] = neighbors === 3 || (alive && neighbors === 2);
    }
  }
  [world, next] = [next, world];
  generation += 1;
  generationNode.textContent = `generation: ${generation}`;
}

function draw() {
  ctx.fillStyle = background;
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const px = x * cellSize;
      const py = y * cellSize;
      ctx.fillStyle = grid;
      ctx.fillRect(px, py, cellSize, cellSize);
      if (!getAlive(world, x, y)) {
        continue;
      }
      ctx.fillStyle = waveColor(x, y, phase);
      ctx.fillRect(px + 1, py + 1, cellSize - 2, cellSize - 2);
    }
  }
}

function waveColor(x, y, phaseValue) {
  let hue = (phaseValue + (x * 0.75 + y * 1.15) / wavelength) % 1;
  if (hue < 0) hue += 1;
  return `hsl(${Math.round(hue * 360)}deg ${saturation}% ${brightness}%)`;
}

function createEmptyWorld() {
  return new Array(width * height).fill(false);
}

function index(x, y) {
  const wrappedX = (x % width + width) % width;
  const wrappedY = (y % height + height) % height;
  return wrappedY * width + wrappedX;
}

function getAlive(cells, x, y) {
  return cells[index(x, y)];
}

function setAlive(cells, x, y, alive) {
  cells[index(x, y)] = alive;
}

function neighborCount(cells, x, y) {
  let count = 0;
  for (let oy = -1; oy <= 1; oy++) {
    for (let ox = -1; ox <= 1; ox++) {
      if (ox === 0 && oy === 0) continue;
      if (getAlive(cells, x + ox, y + oy)) {
        count += 1;
      }
    }
  }
  return count;
}
