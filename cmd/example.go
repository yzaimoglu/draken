package main

import (
	"net/http"
	"time"

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
	apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, req := draken.DrakenHandler(w, r)

		res.Status(http.StatusOK).Json(map[string]any{
			"request_id": req.RequestId(),
			"uptime":     time.Since(d.StartedAt).String(),
		})
	}, draken.TestMiddleware())

	d.Serve()
}
