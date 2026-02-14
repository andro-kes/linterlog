package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules struct {
		CapitalLetter  bool `yaml:"capital_letter"`
		OnlyEnglish    bool `yaml:"only_english"`
		SpecialSymbols bool `yaml:"special_symbols"`
		SensitiveData  bool `yaml:"sensitive_data"`
	} `yaml:"rules"`
}

func ParseConfig() (cfg Config, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	configPath := filepath.Join(wd, "config/config.yml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return
	}

	return
}
