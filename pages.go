package main

import (
	"fmt"
	"io"
	"os"

	"go.yaml.in/yaml/v4"
)

type Pages struct {
	Paths []PageData
}

type PageData struct {
	Label string
	File  string
	Path  string
}

func readPathsFile(file *os.File) ([]Link, error) {
	p, err := parsePathsYaml(file)
	if err != nil {
		return nil, fmt.Errorf("parsing paths file: %w", err)
	}
	names := make([]Link, 0, len(p))
	for _, path := range p {
		names = append(names, NewLink(path.Path, path.Label))
	}
	return names, nil
}

func parsePathsYaml(r io.Reader) ([]PageData, error) {
	p := Pages{}
	y := yaml.NewDecoder(r)
	err := y.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p.Paths, nil
}
