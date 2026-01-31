package handler

import (
	"github.com/a-h/templ"
	"github.com/hnasser-dev/wordle-solver/web/layout"
	"github.com/labstack/echo/v5"
)

type HomeHandler struct{}

func (h HomeHandler) HandleGetHome(c *echo.Context) error {
	return render(c, layout.Home())
}

func render(c *echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}
