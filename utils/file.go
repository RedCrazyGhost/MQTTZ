package utils

import (
	"encoding/json"
	"os"
)

func LoadJSONFile(fileName string, value any) {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(fileBytes, value)
	if err != nil {
		panic(err)
	}
}
