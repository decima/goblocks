package controllers

import (
	"goblocks/app"
	"net/http"
)

type HomeController struct {
	*BaseController
	app *app.App
}

func NewHomeController(app *app.App) *HomeController {
	return &HomeController{
		NewBaseRoute("GET /"),
		app,
	}
}

func (c *HomeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		c.Error(w, "Not found", NotFound)
		return
	}

	c.JSON(w, c.app, Ok)
}
