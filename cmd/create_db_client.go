package main

import (
	"context"
	"fmt"

	"stoke/internal/cfg"
	"stoke/internal/ent"

	_ "github.com/mattn/go-sqlite3"
)

func createDBClient(conf cfg.Config) *ent.Client {
	var dbClient *ent.Client
	var err error

	switch conf.Database.Type {
	case "sqlite", "sqlite3":
		dbClient, err = ent.Open("sqlite3", conf.Database.Sqlite.ConnectionString())
	case "postgres":
		dbClient, err = ent.Open("postgres", conf.Database.Postgres.ConnectionString())
	case "mysql":
		dbClient, err = ent.Open("mysql", conf.Database.Mysql.ConnectionString())
	default:
		err = fmt.Errorf("Unsupported database type: %s", conf.Database.Type)
	}
	if err != nil {
		logger.Error().Err(err).Msg("Could not connect to database")
		panic("Unrecoverable")
	}

	// TODO better migration logic/ check error
	dbClient.Schema.Create(context.Background())
	return dbClient
}
