package settings

import (
	"errors"
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Settings struct {
	FrontendPath string `yaml:"frontend_path,omitempty"`
	ServerPort   int    `yaml:"server_port,omitempty"`
	SiteTitle    string `yaml:"site_title,omitempty"`
}

func ReadSettings(filename string) (Settings, error) {
	s := Settings{"frontend/pages", 9753, "Maze Site"}
	f, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}

	d := yaml.NewDecoder(f)
	err = d.Decode(&s)
	if err != nil {
		return s, err
	}

	fmt.Println(s)
	return s, nil
}
