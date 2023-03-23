package main

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
}

func NewConfig() (config Config, err error) {
	rawToml, err := os.ReadFile("configs/config.toml")
	if err != nil {
		err = fmt.Errorf("failed to read config file: %w", err)
		return
	}

	err = toml.Unmarshal(rawToml, &config)
	if err != nil {
		err = fmt.Errorf("failed to parse config: %w", err)
		return
	}

	return
}
