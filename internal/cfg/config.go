package cfg

import (
	"context"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

type Config struct {
	Server    Server    `json:"server"`
	Database  Database  `json:"database"`
	Logging   Logging   `json:"logging"`
	Tokens    Tokens    `json:"tokens"`
	Users     Users     `json:"users"`
	Telemetry Telemetry `json:"telemetry"`
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
	// We want to do this here to make sure we are populating in the context (though we are using storing the pointer)
	conf.Tokens.ParseDurations()
	return conf
}

func (c *Config) WithContext(ctx context.Context) context.Context {
	confCtx := context.WithValue(ctx, "config", c)
	confCtx = c.Logging.withContext(confCtx)
	confCtx = c.Database.withContext(confCtx)
	confCtx = c.Tokens.withContext(confCtx)
	return c.Users.withContext(confCtx)
}
