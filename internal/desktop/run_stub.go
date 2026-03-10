//go:build !ebitengine

package desktop

import (
	"errors"

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

func Run(Options) error {
	return errors.New("desktop mode requires the ebitengine build tag in this environment: go run -tags ebitengine ./cmd/color-wave-life")
}
