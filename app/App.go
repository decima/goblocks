package app

type App struct {
	Version     string `json:"version"`
	Description string `json:"description"`
}

func NewApp(version string) *App {
	return &App{
		Version:     version,
		Description: "Goblocks Api",
	}
}
