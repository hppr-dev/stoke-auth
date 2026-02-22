package cfg

import (
	"context"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
)

type Config struct {
	Server    Server    `json:"server,omitempty"`
	Database  Database  `json:"database,omitempty"`
	Logging   Logging   `json:"logging,omitempty"`
	Tokens    Tokens    `json:"tokens,omitempty"`
	Users     Users     `json:"users,omitempty"`
	Telemetry Telemetry `json:"telemetry,omitempty"`
	Cluster   Cluster   `json:"cluster,omitempty"`
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
	confCtx := context.WithValue(ctx, configCtxKey, c)
	confCtx = c.Logging.withContext(confCtx)
	confCtx = c.Server.WithContext(confCtx)
	confCtx = c.Database.withContext(confCtx)
	confCtx = c.Cluster.withContext(confCtx)
	confCtx = c.Users.withContext(confCtx)
	return c.Tokens.withContext(confCtx)
}
