package render

import (
	"image/gif"
	"os"
	"path/filepath"
	"testing"
)

func TestExportGIFCreatesExpectedFrames(t *testing.T) {
	dir := t.TempDir()
	outfile := filepath.Join(dir, "wave.gif")

	err := ExportGIF(ExportOptions{
		Width:      12,
		Height:     8,
		Seed:       7,
		Pattern:    "acorn",
		Frames:     6,
		OutputPath: outfile,
		Config: Config{
			CellSize:   4,
			WaveSpeed:  0.05,
			Wavelength: 6,
			Saturation: 0.85,
			Brightness: 0.9,
			Background: DefaultConfig().Background,
			Grid:       DefaultConfig().Grid,
		},
	})
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	info, err := os.Stat(outfile)
	if err != nil {
		t.Fatalf("expected gif to exist: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected gif to be non-empty")
	}

	file, err := os.Open(outfile)
	if err != nil {
		t.Fatalf("open gif: %v", err)
	}
	defer file.Close()

	anim, err := gif.DecodeAll(file)
	if err != nil {
		t.Fatalf("decode gif: %v", err)
	}
	if len(anim.Image) != 6 {
		t.Fatalf("expected 6 frames, got %d", len(anim.Image))
	}
	if anim.Config.Width != 48 || anim.Config.Height != 32 {
		t.Fatalf("unexpected dimensions %dx%d", anim.Config.Width, anim.Config.Height)
	}
}
