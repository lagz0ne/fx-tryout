package main

import (
	"fmt"
	"net/http"

	"go.uber.org/fx"
)

type Config struct {
	Host string
	Port int
}

func resolveConfig() (Config, error) {
	return Config{
		Host: "localhost",
		Port: 3000,
	}, nil
}

func defineViewHostHandler(config Config) HandlerDef {
	return HandlerDef{
		Path: "/host",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("Host: %s", config.Host)))
		},
	}
}

func defineViewPortHandler(config Config) HandlerDef {
	return HandlerDef{
		Path: "/port",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("Port: %d", config.Port)))
		}),
	}
}

var configMod = fx.Module("config",
	fx.Provide(
		resolveConfig,
		fx.Annotate(defineViewHostHandler, fx.ResultTags(`group:"handlers"`)),
		fx.Annotate(defineViewPortHandler, fx.ResultTags(`group:"handlers"`)),
	),
)
