package web

import (
	"net/http"

	"go.uber.org/fx"
)

var Module = fx.Module("web",
	Routes,
	fx.Provide(
		NewHTTPServer,
		fx.Annotate(
			NewRouter,
			fx.ParamTags(`group:"routes"`),
		),
	),
	fx.Invoke(func(*http.Server) {}),
)
