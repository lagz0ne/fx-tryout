package main

import (
	"log"
	"net/http"
	"time"

	"go.uber.org/fx"
)

func registerHandlers(p struct {
	fx.In
	Handlers []HandlerDef `group:"handlers"`
}) {
	for _, h := range p.Handlers {
		log.Printf("Registering handler at %s", h.Path)
		http.Handle(h.Path, h.Handler)
	}
}

type DecoratedHandlers struct {
	fx.Out
	Handlers []HandlerDef `group:"handlers"`
}

func measureHandlerPerformance(p struct {
	fx.In
	Handlers []HandlerDef `group:"handlers"`
}) DecoratedHandlers {
	var handlers []HandlerDef

	for _, h := range p.Handlers {
		path := h.Path
		fn := h.Handler
		handlers = append(handlers, HandlerDef{
			Path: path,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				fn(w, r)
				log.Printf("Handler %s took %s", path, time.Since(start))
			}),
		})
	}

	return DecoratedHandlers{
		Handlers: handlers,
	}
}

var server = fx.Module("server",
	fx.Decorate(measureHandlerPerformance),
	fx.Invoke(
		registerHandlers,
	),
)
