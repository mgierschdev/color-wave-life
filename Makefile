run:
	go run ./cmd/color-wave-life --serve

gif:
	go run ./cmd/color-wave-life --export-gif --pattern glidergun --outfile assets/game-of-life-wave.gif

test:
	go test ./...
