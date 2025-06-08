package utils

import (
	"os"

	json "github.com/bytedance/sonic"
	"gopkg.in/yaml.v3"
)

func LoadJSONFile(fileName string, value any) error {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fileBytes, value)
	if err != nil {
		return err
	}
	return nil
}

func LoadYAMLFile(fileName string, value any) error {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileBytes, value)
	if err != nil {
		return err
	}
	return nil
}
