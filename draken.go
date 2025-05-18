package draken

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Draken struct {
	StartedAt time.Time
	Chi       *chi.Mux
}

func New() *Draken {
	return &Draken{
		StartedAt: time.Now(),
	}
}

func (d *Draken) CreateRouter() {
	d.Chi = chi.NewRouter()
}

func (d *Draken) Get(route string, handler http.HandlerFunc) {
	d.Chi.Get(route, handler)
}

func (d *Draken) Post(route string, handler http.HandlerFunc) {
	d.Chi.Post(route, handler)
}

func (d *Draken) Put(route string, handler http.HandlerFunc) {
	d.Chi.Put(route, handler)
}

func (d *Draken) Patch(route string, handler http.HandlerFunc) {
	d.Chi.Patch(route, handler)
}

func (d *Draken) Delete(route string, handler http.HandlerFunc) {
	d.Chi.Delete(route, handler)
}

func (d *Draken) Serve(addr string) error {
	return http.ListenAndServe(addr, d.Chi)
}

func DrakenHandler(w http.ResponseWriter, r *http.Request) (*Response, *Request) {
	return GetResponse(w), GetRequest(r)
}
