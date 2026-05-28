package main

import "path/filepath"

func GetStaticDir() string {
	return filepath.Join("frontend", "static")
}
