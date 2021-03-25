package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string, out interface{}) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, out)
	if err != nil {
		return err
	}

	return nil
}
