package config

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
	HTTP    HTTPConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
	DSN  string
}

type HTTPConf struct {
	Host, Port string
}

func NewConfig(path string) (config Config, err error) {
	rawToml, err := os.ReadFile(path)
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