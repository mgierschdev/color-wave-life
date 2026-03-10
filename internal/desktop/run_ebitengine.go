//go:build ebitengine

package desktop

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/mgierschdev/color-wave-life/internal/life"
	"github.com/mgierschdev/color-wave-life/internal/render"
)

type Options struct {
	Width         int
	Height        int
	Seed          int64
	Pattern       string
	Density       float64
	SimulationFPS float64
	Config        render.Config
}

type game struct {
	world       *life.World
	image       *ebiten.Image
	buffer      *image.RGBA
	config      render.Config
	width       int
	height      int
	paused      bool
	phase       float64
	simStep     float64
	accumulator float64
	lastTick    time.Time
	currentSeed int64
	fixedSeed   bool
	density     float64
	pattern     string
}

func Run(opts Options) error {
	cfg := opts.Config
	if cfg.CellSize <= 0 {
		cfg = render.DefaultConfig()
	}
	if opts.Width <= 0 {
		opts.Width = 160
	}
	if opts.Height <= 0 {
		opts.Height = 100
	}
	if opts.SimulationFPS <= 0 {
		opts.SimulationFPS = 12
	}
	if opts.Density <= 0 || opts.Density >= 1 {
		opts.Density = life.DefaultDensity
	}

	seed := opts.Seed
	if seed == 0 {
		seed = render.DefaultSeed()
	}

	world := life.NewPatternWorld(opts.Width, opts.Height, opts.Pattern, seed, opts.Density)
	buffer := image.NewRGBA(image.Rect(0, 0, opts.Width*cfg.CellSize, opts.Height*cfg.CellSize))

	g := &game{
		world:       world,
		image:       ebiten.NewImageFromImage(buffer),
		buffer:      buffer,
		config:      cfg,
		width:       opts.Width,
		height:      opts.Height,
		simStep:     1.0 / opts.SimulationFPS,
		lastTick:    time.Now(),
		currentSeed: seed,
		fixedSeed:   opts.Seed != 0,
		density:     opts.Density,
		pattern:     opts.Pattern,
	}
	g.updateTitle()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetTPS(60)
	return ebiten.RunGame(g)
}

func (g *game) Update() error {
	now := time.Now()
	dt := now.Sub(g.lastTick).Seconds()
	g.lastTick = now

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.paused = !g.paused
		g.updateTitle()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.resetWorld()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.resetWithFixedSeed()
	}

	if g.paused {
		return nil
	}

	g.accumulator += dt
	for g.accumulator >= g.simStep {
		g.world.Step()
		g.accumulator -= g.simStep
	}
	g.phase += g.config.WaveSpeed * dt * 60
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	render.DrawFrame(g.buffer, g.world, g.config, g.phase)
	g.image.WritePixels(g.buffer.Pix)
	screen.DrawImage(g.image, nil)
}

func (g *game) Layout(_, _ int) (int, int) {
	return g.buffer.Bounds().Dx(), g.buffer.Bounds().Dy()
}

func (g *game) resetWorld() {
	g.currentSeed = render.DefaultSeed()
	g.world = life.NewPatternWorld(g.width, g.height, g.pattern, g.currentSeed, g.density)
	g.accumulator = 0
	g.phase = 0
	g.updateTitle()
}

func (g *game) resetWithFixedSeed() {
	if g.pattern != "random" {
		g.world = life.NewPatternWorld(g.width, g.height, g.pattern, g.currentSeed, g.density)
		g.accumulator = 0
		g.phase = 0
		g.updateTitle()
		return
	}
	if !g.fixedSeed {
		g.resetWorld()
		return
	}
	rng := rand.New(rand.NewSource(g.currentSeed))
	g.world.Randomize(rng, g.density)
	g.accumulator = 0
	g.phase = 0
	g.updateTitle()
}

func (g *game) updateTitle() {
	status := "running"
	if g.paused {
		status = "paused"
	}
	ebiten.SetWindowTitle(fmt.Sprintf("Color Wave Life | %s | seed=%d | controls: space pause, r randomize, s reseed, q quit", status, g.currentSeed))
}
