package render

import (
	"image"
	"image/gif"
	"os"
	"path/filepath"
	"time"

	"github.com/mgierschdev/color-wave-life/internal/life"
)

type ExportOptions struct {
	Width         int
	Height        int
	Seed          int64
	Pattern       string
	Frames        int
	OutputPath    string
	SimulationFPS float64
	RenderFPS     float64
	Density       float64
	Config        Config
}

func ExportGIF(opts ExportOptions) error {
	cfg := opts.Config
	if cfg.CellSize <= 0 {
		cfg = DefaultConfig()
	}
	if opts.Width <= 0 {
		opts.Width = 120
	}
	if opts.Height <= 0 {
		opts.Height = 80
	}
	if opts.Frames <= 0 {
		opts.Frames = 120
	}
	if opts.OutputPath == "" {
		opts.OutputPath = filepath.Join("assets", "game-of-life-wave.gif")
	}
	if opts.SimulationFPS <= 0 {
		opts.SimulationFPS = 12
	}
	if opts.RenderFPS <= 0 {
		opts.RenderFPS = 20
	}
	if opts.Density <= 0 || opts.Density >= 1 {
		opts.Density = life.DefaultDensity
	}

	world := life.NewPatternWorld(opts.Width, opts.Height, opts.Pattern, opts.Seed, opts.Density)
	canvas := image.NewRGBA(image.Rect(0, 0, opts.Width*cfg.CellSize, opts.Height*cfg.CellSize))
	anim := &gif.GIF{
		Image: make([]*image.Paletted, 0, opts.Frames),
		Delay: make([]int, 0, opts.Frames),
	}

	simStep := 1.0 / opts.SimulationFPS
	renderStep := 1.0 / opts.RenderFPS
	accumulator := 0.0
	phase := 0.0

	for frame := 0; frame < opts.Frames; frame++ {
		accumulator += renderStep
		for accumulator >= simStep {
			world.Step()
			accumulator -= simStep
		}

		DrawFrame(canvas, world, cfg, phase)
		anim.Image = append(anim.Image, Quantize(canvas))
		anim.Delay = append(anim.Delay, int(renderStep*100))
		phase += cfg.WaveSpeed
	}

	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), 0o755); err != nil {
		return err
	}

	file, err := os.Create(opts.OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return gif.EncodeAll(file, anim)
}

func DefaultSeed() int64 {
	return time.Now().UnixNano()
}
