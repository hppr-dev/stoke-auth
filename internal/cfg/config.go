package cfg

import (
	"os"
	"fmt"

	"github.com/ghodss/yaml"
)

type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Logging  Logging  `json:"logging"`
	Tokens   Tokens   `json:"tokens"`
	Users    Users    `json:"users"`
}

func FromFile(filename string) *Config {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file: %s, %v", filename, err))
	}
	conf := &Config{}
	if err = yaml.Unmarshal(content, conf); err != nil {
		panic(fmt.Sprintf("Could not read config yaml: %s, %v", filename, err))
	}
	return conf
}
