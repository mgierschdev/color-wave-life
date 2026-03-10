package life

import "math/rand"

const DefaultDensity = 0.28

type World struct {
	width  int
	height int
	cells  []bool
	next   []bool
}

func NewWorld(width, height int) *World {
	size := width * height
	return &World{
		width:  width,
		height: height,
		cells:  make([]bool, size),
		next:   make([]bool, size),
	}
}

func NewSeededWorld(width, height int, seed int64, density float64) *World {
	if density <= 0 || density >= 1 {
		density = DefaultDensity
	}

	w := NewWorld(width, height)
	w.Randomize(rand.New(rand.NewSource(seed)), density)
	return w
}

func NewPatternWorld(width, height int, pattern string, seed int64, density float64) *World {
	w := NewWorld(width, height)
	if pattern == "" {
		pattern = "glidergun"
	}
	if pattern == "random" {
		w.Randomize(rand.New(rand.NewSource(seed)), density)
		return w
	}
	w.ApplyPattern(pattern)
	return w
}

func (w *World) Width() int {
	return w.width
}

func (w *World) Height() int {
	return w.height
}

func (w *World) Alive(x, y int) bool {
	return w.cells[w.index(x, y)]
}

func (w *World) SetAlive(x, y int, alive bool) {
	w.cells[w.index(x, y)] = alive
}

func (w *World) Cells() []bool {
	return w.cells
}

func (w *World) Clear() {
	for i := range w.cells {
		w.cells[i] = false
	}
}

func (w *World) Randomize(rng *rand.Rand, density float64) {
	if density <= 0 || density >= 1 {
		density = DefaultDensity
	}
	for i := range w.cells {
		w.cells[i] = rng.Float64() < density
	}
}

func (w *World) ApplyPattern(pattern string) {
	w.Clear()
	cx := w.width / 2
	cy := w.height / 2

	var points [][2]int
	switch pattern {
	case "spacefiller":
		cx = w.width/2 - 24
		cy = w.height/2 - 13
		points = [][2]int{
			{-4, -11}, {-4, -10}, {-4, -9}, {0, -11}, {0, -10}, {0, -9},
			{-5, -10}, {-1, -10},
			{-24, -9}, {-23, -9}, {-22, -9}, {-21, -9},
			{-24, -8}, {-21, -8},
			{-24, -3}, {-21, -3},
			{-18, -3},
			{-17, -2}, {-15, -2},
			{-24, -1}, {-21, -1}, {-19, -1}, {-18, -1},
			{-17, -1}, {-16, -1}, {-15, -1}, {-14, -1}, {-13, -1},
			{-21, 0}, {-20, 0},
			{-12, 0}, {-11, 0},
			{-24, 1}, {-23, 1}, {-22, 1}, {-21, 1},
			{-10, 1}, {-9, 1}, {-8, 1}, {-7, 1},
			{-18, 2}, {-12, 2}, {-11, 2},
			{-5, 2}, {-3, 2}, {-2, 2}, {0, 2},
			{-24, 3}, {-12, 3}, {-8, 3}, {-5, 3}, {-4, 3}, {-3, 3}, {-1, 3},
			{-24, 4}, {-22, 4}, {-20, 4}, {-19, 4}, {-18, 4}, {-13, 4}, {-7, 4},
			{-15, 5}, {-11, 5}, {-10, 5}, {-9, 5}, {-8, 5}, {-5, 5}, {-3, 5},
			{-18, 6}, {-17, 6}, {-15, 6}, {-14, 6}, {-10, 6}, {-7, 6}, {-5, 6}, {-4, 6},
			{-3, 6}, {-2, 6}, {-1, 6}, {0, 6},
			{-18, 7}, {-16, 7}, {-15, 7}, {-14, 7}, {-10, 7}, {-9, 7}, {-8, 7}, {-5, 7},
			{-3, 7}, {-1, 7},
			{-19, 8}, {-18, 8}, {-17, 8}, {-15, 8}, {-10, 8}, {-5, 8}, {-3, 8},
			{-2, 8}, {-1, 8}, {0, 8},
			{-21, 9}, {-20, 9}, {-19, 9}, {-18, 9}, {-14, 9}, {-13, 9}, {-12, 9},
			{-11, 9}, {-10, 9},
			{-24, 10}, {-23, 10}, {-22, 10}, {-21, 10},
			{-15, 11}, {-14, 11}, {-13, 11},
			{-18, 12},
			{-24, 13}, {-23, 13}, {-22, 13}, {-21, 13},
		}
	case "glider":
		points = [][2]int{
			{0, -1},
			{1, 0},
			{-1, 1}, {0, 1}, {1, 1},
		}
	case "glidergun":
		cx = w.width/2 - 18
		cy = w.height/2 - 4
		points = [][2]int{
			{0, 4}, {0, 5}, {1, 4}, {1, 5},
			{10, 4}, {10, 5}, {10, 6},
			{11, 3}, {11, 7},
			{12, 2}, {12, 8},
			{13, 2}, {13, 8},
			{14, 5},
			{15, 3}, {15, 7},
			{16, 4}, {16, 5}, {16, 6},
			{17, 5},
			{20, 2}, {20, 3}, {20, 4},
			{21, 2}, {21, 3}, {21, 4},
			{22, 1}, {22, 5},
			{24, 0}, {24, 1}, {24, 5}, {24, 6},
			{34, 2}, {34, 3}, {35, 2}, {35, 3},
		}
	case "lwss":
		points = [][2]int{
			{-2, -1}, {-1, -1}, {0, -1}, {1, -1},
			{-3, 0}, {1, 0},
			{1, 1},
			{-3, 2}, {0, 2},
		}
	case "pulsar":
		points = [][2]int{
			{-4, -6}, {-3, -6}, {-2, -6}, {2, -6}, {3, -6}, {4, -6},
			{-6, -4}, {-1, -4}, {1, -4}, {6, -4},
			{-6, -3}, {-1, -3}, {1, -3}, {6, -3},
			{-6, -2}, {-1, -2}, {1, -2}, {6, -2},
			{-4, -1}, {-3, -1}, {-2, -1}, {2, -1}, {3, -1}, {4, -1},
			{-4, 1}, {-3, 1}, {-2, 1}, {2, 1}, {3, 1}, {4, 1},
			{-6, 2}, {-1, 2}, {1, 2}, {6, 2},
			{-6, 3}, {-1, 3}, {1, 3}, {6, 3},
			{-6, 4}, {-1, 4}, {1, 4}, {6, 4},
			{-4, 6}, {-3, 6}, {-2, 6}, {2, 6}, {3, 6}, {4, 6},
		}
	case "diehard":
		points = [][2]int{
			{-3, 0},
			{-2, 0},
			{-2, 1},
			{2, 1},
			{3, -1}, {3, 1},
			{4, 1},
		}
	case "rpentomino":
		points = [][2]int{
			{0, -1}, {1, -1},
			{-1, 0}, {0, 0},
			{0, 1},
		}
	case "switchengine":
		points = [][2]int{
			{-4, -2}, {-3, -2},
			{-2, -1}, {-4, -1},
			{-4, 0},
			{-2, 1}, {0, 1}, {1, 1},
			{2, 0},
			{0, -2},
		}
	default:
		points = [][2]int{
			{-3, 0},
			{-1, 1},
			{-3, 2}, {-2, 2}, {1, 2},
			{-3, 3},
		}
	}

	for _, p := range points {
		w.SetAlive(cx+p[0], cy+p[1], true)
	}
}

func (w *World) Step() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			neighbors := w.neighborCount(x, y)
			alive := w.Alive(x, y)
			w.next[w.index(x, y)] = neighbors == 3 || (alive && neighbors == 2)
		}
	}

	w.cells, w.next = w.next, w.cells
}

func (w *World) index(x, y int) int {
	x = ((x % w.width) + w.width) % w.width
	y = ((y % w.height) + w.height) % w.height
	return y*w.width + x
}

func (w *World) neighborCount(x, y int) int {
	count := 0
	for oy := -1; oy <= 1; oy++ {
		for ox := -1; ox <= 1; ox++ {
			if ox == 0 && oy == 0 {
				continue
			}
			if w.Alive(x+ox, y+oy) {
				count++
			}
		}
	}
	return count
}
