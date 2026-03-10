package web

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

//go:embed static/*
var staticFS embed.FS

type Options struct {
	OpenBrowser bool
}

func Serve(opts Options) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	url := "http://" + listener.Addr().String() + "/static/"
	fmt.Println("Open:", url)

	if opts.OpenBrowser && runtime.GOOS == "darwin" {
		_ = exec.Command("open", url).Start()
	}

	return http.Serve(listener, mux)
}
