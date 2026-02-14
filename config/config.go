package config

import (
	"os"

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

func ParseConfig() (cfg *Config, err error) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
