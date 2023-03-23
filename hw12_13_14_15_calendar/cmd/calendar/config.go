package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Storage StorageConf `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageType string

const (
	StorageTypeMemory StorageType = "memory"
	StorageTypeSql    StorageType = "sql"
)

type StorageConf struct {
	Type StorageType `yaml:"type"`
}

func NewConfig() (config Config, err error) {
	rawYaml, err := os.ReadFile("config.yml")
	if err != nil {
		err = fmt.Errorf("failed to read config file: %w", err)
		return
	}

	err = yaml.Unmarshal(rawYaml, &config)
	if err != nil {
		err = fmt.Errorf("failed to parse config: %w", err)
		return
	}

	return
}
