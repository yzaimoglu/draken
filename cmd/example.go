package main

import (
	"net/http"
	"time"

	"github.com/yzaimoglu/draken"
)

func main() {
	d := draken.New()
	d.CreateRouter()
	d.EssentialMiddlewares()
	d.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res, req := draken.DrakenHandler(w, r)

		res.Status(http.StatusOK).Json(map[string]any{
			"request_id": req.RequestId(),
			"uptime":     time.Since(d.StartedAt).String(),
		})
	})
	d.Serve(":3000")
}
