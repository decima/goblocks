package web

import (
	"bytes"
	"fmt"
	"goblocks/app/web/controllers"
	"log/slog"
	"net/http"
	"time"

	"go.uber.org/fx"
)

var Routes = fx.Module("router",
	AsRoute(controllers.NewHomeController),
	AsRoutes(
		controllers.NewGetBlockController,
		controllers.NewWriteBlockController,
		controllers.NewDeleteBlockController),
	fx.Provide(),
)

func NewRouter(routes []Route, logger *slog.Logger) *Router {
	router := http.NewServeMux()
	for _, route := range routes {
		router.Handle(route.Pattern(), route)
	}

	return &Router{logger, router}
}

type Route interface {
	http.Handler
	// Pattern returns the mux pattern of the route with the go 1.22+ syntax
	Pattern() string
}

func AsRoute(f any) fx.Option {
	return fx.Module(fmt.Sprintf("controller"), fx.Provide(
		fx.Annotate(
			f,
			fx.As(new(Route)),
			fx.ResultTags(`group:"routes"`),
		),
	),
	)
}

func AsRoutes(f ...any) fx.Option {
	options := []fx.Option{}
	for _, v := range f {
		options = append(options, AsRoute(v))
	}
	return fx.Options(options...)
}

type Router struct {
	logger   *slog.Logger
	serveMux *http.ServeMux
}

func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	buf := &responseBuffer{
		header:     make(http.Header),
		buffer:     new(bytes.Buffer),
		statusCode: 0,
	}
	start := time.Now()

	shouldPrintOutput := true

	defer func() {
		if err := recover(); err != nil {
			if r.logger != nil {
				r.logger.ErrorContext(req.Context(), "panic", "recover", err)
			}

			controllers.Error(w, "Something went wrong")
			shouldPrintOutput = false
		}

		if shouldPrintOutput && (buf.statusCode == 0 || buf.statusCode == http.StatusNotFound) {
			if buf.buffer.Len() == 0 {
				controllers.Error(w, "Not found", controllers.NotFound)
				shouldPrintOutput = false
			}
		}

		elapsed := float64(time.Since(start).Nanoseconds()) / 100000000
		r.logger.Info(fmt.Sprintf(
			"%v %v %v %v",
			req.Method,
			req.URL.Path,
			buf.statusCode,
			elapsed,
		),
			"http.method", req.Method,
			"http.path", req.URL.Path,
			"http.pattern", req.Pattern,
			"http.user-agent", req.UserAgent(),
			"http.status", buf.statusCode,
			"http.duration", elapsed,
		)

		if !shouldPrintOutput {
			return
		}

		for k, vv := range buf.header {
			w.Header()[k] = vv
		}

		w.WriteHeader(buf.statusCode)
		buf.buffer.WriteTo(w)
	}()

	r.serveMux.ServeHTTP(buf, req)
}

type responseBuffer struct {
	header      http.Header
	buffer      *bytes.Buffer
	statusCode  int
	wroteHeader bool
}

func (rb *responseBuffer) Header() http.Header {
	return rb.header
}

func (rb *responseBuffer) Write(b []byte) (int, error) {
	return rb.buffer.Write(b)
}

func (rb *responseBuffer) WriteHeader(statusCode int) {
	if !rb.wroteHeader {
		rb.statusCode = statusCode
		rb.wroteHeader = true
	}
}
