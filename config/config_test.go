package config

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		fmt.Println("fail to read file")
		t.Fail()
	}

	if len(data) == 0 {
		fmt.Println("file is empty")
		t.Fail()
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
}