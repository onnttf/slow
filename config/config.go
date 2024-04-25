package config

import (
	"fmt"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

type ConfigType struct {
	Database struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
}

var (
	once   sync.Once
	Config *ConfigType
)

const filePath = "config/config.json"

func Initialize() error {
	var err error
	once.Do(func() {
		var content []byte
		content, err = os.ReadFile(filePath)
		if err != nil {
			err = fmt.Errorf("failed to read file: %s", err.Error())
			return
		}
		if err = jsoniter.Unmarshal(content, &Config); err != nil {
			err = fmt.Errorf("failed to unmarshal JSON: %s", err.Error())
			return
		}
	})
	return err // Return the captured error value
}
