package main

import (
	"fmt"
	"testing"

	"go.uber.org/fx"
)

func TestConfig(t *testing.T) {

	t.Run("test config get", func(t *testing.T) {
		var config Config

		fx.New(
			configMod,
			fx.Populate(&config),
		)

		if config.Host != "localhost" {
			t.FailNow()
		}
	})

	t.Run("test replace", func(t *testing.T) {
		var config Config

		fx.New(
			configMod,
			fx.Replace(Config{
				Host: "test",
				Port: 6000,
			}),
			fx.Populate(&config),
		)

		fmt.Printf("Detail %v+", config)

		if config.Host != "test" {
			t.FailNow()
		}

	})

}
