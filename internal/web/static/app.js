const canvas = document.getElementById("life");
const ctx = canvas.getContext("2d");
const toggle = document.getElementById("toggle");
const reset = document.getElementById("reset");
const slower = document.getElementById("slower");
const faster = document.getElementById("faster");
const record = document.getElementById("record");
const patternSelect = document.getElementById("pattern");
const statusNode = document.getElementById("status");
const patternLabelNode = document.getElementById("pattern-label");
const generationNode = document.getElementById("generation");
const speedNode = document.getElementById("speed");

const width = 240;
const height = 160;
const cellSize = 4;
const minSimulationFPS = 1;
const maxSimulationFPS = 100;
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
let currentPattern = "glidergun";

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

slower.addEventListener("click", () => {
  simulationFPS = Math.max(minSimulationFPS, simulationFPS - 1);
  simulationAccumulator = 0;
  updateSpeedLabel();
});

faster.addEventListener("click", () => {
  simulationFPS = Math.min(maxSimulationFPS, simulationFPS + 1);
  simulationAccumulator = 0;
  updateSpeedLabel();
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

function updateSpeedLabel() {
  speedNode.textContent = `speed: ${simulationFPS.toFixed(1)} gen/s`;
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
