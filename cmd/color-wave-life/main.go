package main

import (
	"flag"
	"log"
	"strings"

	"github.com/mgierschdev/color-wave-life/internal/desktop"
	"github.com/mgierschdev/color-wave-life/internal/life"
	"github.com/mgierschdev/color-wave-life/internal/render"
	"github.com/mgierschdev/color-wave-life/internal/web"
)

func main() {
	width := flag.Int("width", 120, "board width in cells")
	height := flag.Int("height", 80, "board height in cells")
	cellSize := flag.Int("cell-size", 8, "render size of each cell in pixels")
	seed := flag.Int64("seed", 0, "deterministic seed; 0 picks a fresh random seed")
	speed := flag.Float64("speed", 0.045, "wave travel speed per frame")
	wavelength := flag.Float64("wavelength", 11, "wave length across the grid")
	saturation := flag.Float64("saturation", 0.85, "cell color saturation from 0 to 1")
	brightness := flag.Float64("brightness", 0.92, "cell brightness from 0 to 1")
	frames := flag.Int("frames", 140, "number of frames to export when --export-gif is set")
	outfile := flag.String("outfile", "assets/game-of-life-wave.gif", "gif output path")
	exportGIF := flag.Bool("export-gif", false, "render a deterministic gif instead of opening the desktop window")
	serve := flag.Bool("serve", false, "open a browser preview with built-in recording controls")
	openBrowser := flag.Bool("open-browser", true, "open the browser automatically in serve mode")
	pattern := flag.String("pattern", "glidergun", "starting pattern: glidergun, rpentomino, acorn, or random")
	density := flag.Float64("density", life.DefaultDensity, "starting density of alive cells from 0 to 1")
	simulationFPS := flag.Float64("simulation-fps", 12, "simulation updates per second")
	renderFPS := flag.Float64("render-fps", 20, "frame rate used by gif capture")
	flag.Parse()

	cfg := render.DefaultConfig()
	cfg.CellSize = *cellSize
	cfg.WaveSpeed = *speed
	cfg.Wavelength = *wavelength
	cfg.Saturation = *saturation
	cfg.Brightness = *brightness

	if *exportGIF {
		exportSeed := *seed
		if exportSeed == 0 {
			exportSeed = 20260310
		}
		err := render.ExportGIF(render.ExportOptions{
			Width:         *width,
			Height:        *height,
			Seed:          exportSeed,
			Pattern:       strings.ToLower(*pattern),
			Frames:        *frames,
			OutputPath:    *outfile,
			SimulationFPS: *simulationFPS,
			RenderFPS:     *renderFPS,
			Density:       *density,
			Config:        cfg,
		})
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *serve {
		if err := web.Serve(web.Options{OpenBrowser: *openBrowser}); err != nil {
			log.Fatal(err)
		}
		return
	}

	err := desktop.Run(desktop.Options{
		Width:         *width,
		Height:        *height,
		Seed:          *seed,
		Pattern:       strings.ToLower(*pattern),
		Density:       *density,
		SimulationFPS: *simulationFPS,
		Config:        cfg,
	})
	if err != nil {
		if err := web.Serve(web.Options{OpenBrowser: *openBrowser}); err != nil {
			log.Fatal(err)
		}
	}
}
