# color-wave-life

![Color Wave Life](assets/game-of-life-wave.gif)

`color-wave-life` is a Go implementation of Conway's Game of Life with a traveling color field layered over each living cell. The browser preview now includes a `spacefiller` pattern for sustained expansion instead of quick stabilization.

## Features

- Desktop app built with Ebitengine.
- Browser preview with built-in recording to WebM.
- Deterministic seeded runs for reproducible patterns.
- Smooth color-wave rendering driven by time and cell position.
- GIF export mode for generating README-ready captures.
- Toroidal grid wrapping so patterns flow across edges.
- Default `spacefiller` pattern for sustained expansion, plus `glidergun`, `switchengine`, `pulsar`, `rpentomino`, `acorn`, `diehard`, `lwss`, and `glider`.

## Run locally

```bash
go run ./cmd/color-wave-life --serve
```

That opens a local browser page automatically on macOS.

Useful flags:

```bash
go run -tags ebitengine ./cmd/color-wave-life \
  --width 140 \
  --height 90 \
  --cell-size 7 \
  --speed 0.05 \
  --wavelength 12 \
  --pattern glidergun
```

## Controls

- `Space`: pause or resume
- `R`: randomize with a fresh seed
- `S`: reseed deterministically when `--seed` is provided
- `Q` or `Esc`: quit

## Generate the GIF

```bash
go run ./cmd/color-wave-life \
  --export-gif \
  --pattern glidergun \
  --width 120 \
  --height 80 \
  --cell-size 6 \
  --frames 140 \
  --outfile assets/game-of-life-wave.gif
```

The exported animation is deterministic when a seed is provided, which makes it suitable for the README and release assets.

## Browser preview

```bash
go run ./cmd/color-wave-life --serve
```

The browser preview includes:

- live animation from the very first cells
- pattern switching between `spacefiller`, `glidergun`, `switchengine`, `pulsar`, `rpentomino`, `acorn`, `diehard`, `lwss`, and `glider`
- direct `Gen/s` numeric input with no fixed upper cap
- `Pause` and `Reset` controls
- `Record 12s` to download a WebM capture directly from the canvas

## Development

```bash
go test ./...
```

If your local macOS graphics toolchain is healthy, run the desktop window with the `ebitengine` build tag. Otherwise, the browser preview and GIF export path work without the native windowing layer.

## Color-wave idea

The visual layer uses a hue field based on both time and grid position:

`hue = time_phase + spatial_phase(x, y)`

That phase relationship makes color drift diagonally through the board, so the simulation reads like a flowing current rather than a static rainbow mask.
