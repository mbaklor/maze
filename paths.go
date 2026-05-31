package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

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
