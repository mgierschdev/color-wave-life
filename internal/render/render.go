package render

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math"

	"github.com/mgierschdev/color-wave-life/internal/life"
)

type Config struct {
	CellSize   int
	WaveSpeed  float64
	Wavelength float64
	Saturation float64
	Brightness float64
	Background color.RGBA
	Grid       color.RGBA
}

func DefaultConfig() Config {
	return Config{
		CellSize:   8,
		WaveSpeed:  0.045,
		Wavelength: 11,
		Saturation: 0.85,
		Brightness: 0.92,
		Background: color.RGBA{R: 8, G: 10, B: 18, A: 255},
		Grid:       color.RGBA{R: 20, G: 24, B: 36, A: 255},
	}
}

func DrawFrame(dst *image.RGBA, world *life.World, cfg Config, phase float64) {
	bounds := dst.Bounds()
	fill(dst, bounds, cfg.Background)

	cellSize := maxInt(cfg.CellSize, 1)
	for y := 0; y < world.Height(); y++ {
		for x := 0; x < world.Width(); x++ {
			rect := image.Rect(
				x*cellSize,
				y*cellSize,
				(x+1)*cellSize,
				(y+1)*cellSize,
			)

			draw.Draw(dst, rect, &image.Uniform{C: cfg.Grid}, image.Point{}, draw.Src)
			if !world.Alive(x, y) {
				continue
			}

			inset := maxInt(cellSize/8, 1)
			if cellSize <= 3 {
				inset = 0
			}
			inner := rect.Inset(inset)
			if inner.Empty() {
				inner = rect
			}
			draw.Draw(dst, inner, &image.Uniform{C: WaveColor(x, y, phase, cfg)}, image.Point{}, draw.Src)
		}
	}
}

func WaveColor(x, y int, phase float64, cfg Config) color.RGBA {
	wavelength := cfg.Wavelength
	if wavelength <= 0 {
		wavelength = 1
	}
	spatial := (float64(x)*0.75 + float64(y)*1.15) / wavelength
	hue := math.Mod(phase+spatial, 1.0)
	if hue < 0 {
		hue += 1.0
	}
	return hsvToRGBA(hue, clamp(cfg.Saturation, 0, 1), clamp(cfg.Brightness, 0, 1))
}

func Quantize(src image.Image) *image.Paletted {
	paletted := image.NewPaletted(src.Bounds(), palette.Plan9)
	draw.FloydSteinberg.Draw(paletted, src.Bounds(), src, image.Point{})
	return paletted
}

func hsvToRGBA(h, s, v float64) color.RGBA {
	if s == 0 {
		g := uint8(math.Round(v * 255))
		return color.RGBA{R: g, G: g, B: g, A: 255}
	}

	h = math.Mod(h, 1) * 6
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	var r, g, b float64
	switch int(i) {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}

	return color.RGBA{
		R: uint8(math.Round(r * 255)),
		G: uint8(math.Round(g * 255)),
		B: uint8(math.Round(b * 255)),
		A: 255,
	}
}

func fill(dst *image.RGBA, bounds image.Rectangle, c color.RGBA) {
	draw.Draw(dst, bounds, &image.Uniform{C: c}, image.Point{}, draw.Src)
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
