package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DisableLowercase       bool   `yaml:"disable-lowercase"`
	DisableEnglishOnly     bool   `yaml:"disable-english-only"`
	DisableNoSpecialChars  bool   `yaml:"disable-no-special-chars"`
	DisableNoSensitiveData bool   `yaml:"disable-no-sensitive-data"`
	ExtraSensitivePatterns string `yaml:"extra-sensitive-patterns"`
}

func LoadConfig() (*Config, error) {
	// Пытаемся найти путь к конфигу относительно файла config.go
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	absConfigPath := filepath.Join(projectRoot, "configs/config.yaml")

	var cfg Config
	if err := cleanenv.ReadConfig(absConfigPath, &cfg); err != nil {
		// Если не нашли по абсолютному пути, пробуем относительный
		if errRel := cleanenv.ReadConfig("configs/config.yaml", &cfg); errRel != nil {
			return nil, fmt.Errorf("failed to read config (tried %s and relative): %w", absConfigPath, errRel)
		}
	}

	return &cfg, nil
}
