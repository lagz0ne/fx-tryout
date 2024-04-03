package main

import (
	"context"
	"net/http"

	"go.uber.org/fx"
)

type HandlerDef struct {
	Path    string
	Handler http.HandlerFunc
}

func main() {
	app := fx.New(
		db,
		configMod,
		person,
		server,
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						http.ListenAndServe(":3000", nil)
					}()

					return nil
				},
			})
		}),
	)

	app.Run()

}
