package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yzaimoglu/draken"
)

func main() {
	d, err := draken.New()
	if err != nil {
		panic(err)
	}
	d.CreateRouter()
	d.Router.EssentialMiddlewares()

	apiRouter := d.Router.CreateSubrouter("/api/v1")
	apiRouter.Get("", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]any{
			"request_id": ctx.Get(string(draken.ContextKeyRequestId)),
			"uptime":     time.Since(d.StartedAt).String(),
		})
	})

	d.Serve()
}
