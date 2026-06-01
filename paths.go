package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Link struct {
	Path string
	Name string
}

func NewLink(path, name string) Link {
	return Link{Path: path, Name: name}
}

func GetStaticDir() string {
	return filepath.Join("frontend", "static")
}

func ParseFileFromUrl(url string) (string, error) {
	path := filepath.Join("frontend", "pages", url)
	isDir, err := pathIsDir(path)
	if err != nil {
		return "", err
	}
	if isDir {
		return filepath.Join(path, "index.md"), nil
	}
	return path + ".md", nil
}

func BasePathFromUrl(url string) string {
	if url == "/" {
		return url
	}
	s := strings.Split(url, "/")
	return "/" + s[1]
}

func pathIsDir(path string) (bool, error) {
	st, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return st.IsDir(), nil
}

func findBase() ([]Link, error) {
	root := filepath.Join("frontend", "pages")
	f, err := os.Open(filepath.Join(root, "pages.yml"))
	if errors.Is(err, os.ErrNotExist) {
		return readRootDir(root)
	}
	if err != nil {
		return nil, fmt.Errorf("opening paths yaml: %w", err)
	}
	return readPathsFile(f)
}

func readRootDir(root string) ([]Link, error) {
	dir, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	names := make([]Link, 0, len(dir))
	names = append(names, NewLink("/", "Home"))
	for _, f := range dir {
		if f.IsDir() {
			names = append(names, NewLink("/"+f.Name(), f.Name()))
		}
	}
	return names, nil
}
