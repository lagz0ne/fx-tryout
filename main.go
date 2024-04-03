package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"
)

type HandlerDef struct {
	Path    string
	Handler http.HandlerFunc
}

func main() {
	fx.New(
		db,
		configMod,
		person,
		fx.Invoke(func(lifecycle fx.Lifecycle, p struct {
			fx.In
			Config   Config
			Handlers []HandlerDef `group:"handlers"`
		}) {
			for _, h := range p.Handlers {
				log.Printf("Registering handler at %s", h.Path)
				http.Handle(h.Path, h.Handler)
			}

			lifecycle.Append(fx.Hook{
				OnStart: func(context.Context) error {
					log.Printf("Server started at %s:%d", p.Config.Host, p.Config.Port)
					return nil
				},
			})
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", p.Config.Port), nil))
		}),
	).Run()
}
