package main

import (
	"goblocks/app"
	"goblocks/app/config"
	"goblocks/app/services"
	"goblocks/app/web"
	"goblocks/libraries/utils/prettylog"
	"log/slog"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var Version string = "0.0"
var Debug = os.Getenv("DEBUG") == "1"

func main() {
	log := NewLogger()

	fx.New(
		fx.Supply(app.NewApp(Version)),
		fx.Supply(log),
		fx.Provide(config.NewConfig),
		web.Module,
		services.Module,
		fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: log}
		}),
	).Run()
}

func NewLogger() *slog.Logger {
	if Debug {
		opts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		return slog.New(prettylog.NewHandler(opts))

	}
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

//mux := http.NewServeMux()
//mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
//	io.WriteString(w, fmt.Sprintf("Hello %s", r.PathValue("name")))
//})
//
//mux.HandleFunc("GET /demo/{any}", func(w http.ResponseWriter, r *http.Request) {
//	io.WriteString(w, fmt.Sprintf("Hello Demo"))
//})
//
//apiServer := &http.Server{
//	Addr:           ":8000", // just informational since we created the listener ourselves
//	Handler:        mux,
//	ReadTimeout:    10 * time.Second,
//	WriteTimeout:   10 * time.Second,
//	MaxHeaderBytes: 16384,
//}
//log.Fatal(apiServer.ListenAndServe())
