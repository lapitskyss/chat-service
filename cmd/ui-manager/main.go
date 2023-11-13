package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

const addr = ":3001"

//go:embed static
var staticFiles embed.FS

func main() {
	log.Printf("Look at http://localhost%v/", addr)
	fSys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(addr, http.FileServer(http.FS(fSys))); err != nil { //nolint:gosec
		log.Fatal(err)
	}
}
