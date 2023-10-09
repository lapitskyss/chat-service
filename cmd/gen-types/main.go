package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("invalid args count: %d", len(os.Args)-1)
	}

	pkg, types, out := os.Args[1], strings.Split(os.Args[2], ","), os.Args[3]
	if err := run(pkg, types, out); err != nil {
		log.Fatal(err)
	}

	p, _ := os.Getwd()
	fmt.Printf("%v generated\n", filepath.Join(p, out)) //nolint:forbidigo
}

//go:embed templates/types.go.tpl
var templates embed.FS

func run(pkg string, types []string, outFile string) error {
	tmpl := template.Must(template.ParseFS(templates, "templates/types.go.tpl"))

	tplContext := map[string]interface{}{
		"packageName": pkg,
		"types":       types,
	}
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, tplContext); err != nil {
		return fmt.Errorf("execute template: %v", err)
	}

	const perm = 0o644
	if err := os.WriteFile(outFile, buf.Bytes(), perm); err != nil {
		return fmt.Errorf("cannot write result: %v", err)
	}

	return nil
}
