package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
)

type DbConfig struct {
	Url string
}

type dbContainer struct {
	db *sqlx.DB
}

type DbConnection interface {
	sqlx.DB
}

func provideConfigOnDev(lc fx.Lifecycle) *DbConfig {
	ctx := context.Background()
	dev := os.Getenv("APP_ENV")
	if dev == "dev" {
		log.Printf("Running in dev mode")
		container := testcontainers.WithImage("docker.io/postgres:15.2-alpine")

		postgresContainer, e := postgres.RunContainer(ctx,
			container,
			postgres.WithDatabase("test"),
			postgres.WithUsername("postgres"),
			postgres.WithPassword("password"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(5*time.Second)),
		)

		lc.Append(fx.Hook{
			OnStop: func(context.Context) error {
				return postgresContainer.Terminate(ctx)
			},
		})

		if e != nil {
			panic(e)
		}

		log.Printf("Container Detail: %v+", postgresContainer)

		url, e := postgresContainer.ConnectionString(ctx, "sslmode=disable", "dbname=test")
		if e != nil {
			panic(e)
		}
		log.Printf("Running in dev mode")

		log.Printf("Connection string: %s", url)
		db := sqlx.MustConnect("postgres", url)
		defer db.Close()

		log.Printf("Creating schema")
		_, e = db.Exec(schema)
		if e != nil {
			panic(e)
		}

		log.Printf("Successfully created schema")

		return &DbConfig{
			Url: url,
		}
	}

	return nil
}

func RetrieveConfig(config *DbConfig) (DbConfig, error) {
	if config != nil {
		return *config, nil
	}

	return DbConfig{
		Url: "user=postgres password=password dbname=postgres sslmode=disable",
	}, nil
}

func ProvideConnection(config DbConfig) (*dbContainer, error) {
	db := sqlx.MustConnect("postgres", config.Url)

	return &dbContainer{
		db,
	}, nil
}

var db = fx.Module("db",
	fx.Provide(
		provideConfigOnDev,
		RetrieveConfig,
		ProvideConnection,
	),
)
