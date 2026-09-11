package settings

import (
	"errors"
	"os"

	"github.com/mbaklor/go-env"
	"go.yaml.in/yaml/v4"
)

type Settings struct {
	ConfigFile   string `yaml:"-" env:"CONFIG_FILE"`
	FrontendPath string `yaml:"frontend_path,omitempty" env:"FRONTEND_PATH"`
	ServerPort   int    `yaml:"server_port,omitempty" env:"SERVER_PORT"`
	SiteTitle    string `yaml:"site_title,omitempty" env:"SITE_TITLE"`
}

func ReadSettingsFile(filename string) (Settings, error) {
	var s Settings
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

	return s, nil
}

func ReadEnvVars() (Settings, error) {
	var c Settings
	err := env.Load(&c)
	if err != nil {
		return c, err
	}
	return c, nil
}
