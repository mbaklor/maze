package paths

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

func FrontendPath(p ...string) string {
	elems := make([]string, 0, len(p)+1)
	elems = append(elems, "frontend")
	elems = append(elems, p...)
	return filepath.Join(elems...)
}

func ParseFileFromUrl(url, root string) (string, error) {
	path := filepath.Join(root, url)
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

func FindBase(root string) ([]Link, error) {
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
