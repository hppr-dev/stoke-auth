package cfg

import (
	"context"
	"fmt"
	"stoke/internal/ent"

	"github.com/rs/zerolog"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	// One Of sqlite3, postgres, mysql
	Type     string   `json:"type"`
	// sqlite config
	Sqlite   Sqlite  `json:"sqlite"`
	// postgres config
	Postgres Postgres `json:"postgres"`
	// mysql config
	Mysql    Mysql    `json:"mysql"`
}

type Sqlite struct {
	File string  `json:"file"`
	Flags string `json:"flags"`
}

func (s Sqlite) ConnectionString() string {
	return fmt.Sprintf("file:%s?%s", s.File, s.Flags)
}

type Postgres struct {
	Host     	string `json:"host"`
	Port     	int    `json:"port"`
	Database 	string `json:"database"`
	User     	string `json:"user"`
	Password 	string `json:"password"`
	// String to append to the end of the connection string
	ExtraArgs string `json:"extra_args"`
}

func (p Postgres) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s %s", p.Host, p.Port, p.User, p.Database, p.Password, p.ExtraArgs)
}

type Mysql struct {
	Host     	string `json:"host"`
	Port     	int    `json:"port"`
	Database 	string `json:"database"`
	User     	string `json:"user"`
	Password 	string `json:"password"`
	Flags     string `json:"flag"`
}

func (m Mysql) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", m.User, m.Password, m.Host, m.Port, m.Database, m.Flags)
}

func (d Database) withContext(ctx context.Context) context.Context {
	var dbClient *ent.Client
	var err error

	switch d.Type {
	case "sqlite", "sqlite3":
		dbClient, err = ent.Open("sqlite3", d.Sqlite.ConnectionString())
	case "postgres":
		dbClient, err = ent.Open("postgres", d.Postgres.ConnectionString())
	case "mysql":
		dbClient, err = ent.Open("mysql", d.Mysql.ConnectionString())
	default:
		err = fmt.Errorf("Unsupported database type: %s", d.Type)
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Could not connect to database")
		panic("Unrecoverable")
	}

	// TODO better migration logic/ check error
	dbClient.Schema.Create(context.Background())

	return ent.NewContext(ctx, dbClient)
}
