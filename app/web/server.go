package web

import (
	"context"
	"errors"
	"fmt"
	"goblocks/app/config"
	"log/slog"
	"net"
	"net/http"

	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, conf *config.Config, router *Router, log *slog.Logger) *http.Server {
	srv := &http.Server{
		Addr:    conf.HttpHostAndPort(),
		Handler: router,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info(fmt.Sprintf("Listening on %s", srv.Addr))
			go func() {
				err := srv.Serve(ln)
				if err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}
