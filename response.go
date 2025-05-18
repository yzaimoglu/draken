package draken

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	statusCode int
}

func GetResponse(w http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (r *Response) Status(code int) *Response {
	r.statusCode = code
	return r
}

func (r *Response) SetHeader(header string, value string) *Response {
	r.Header().Set(header, value)
	return r
}

func (r *Response) SetContentType(t string) *Response {
	r.SetHeader("Content-Type", t)
	return r
}

func (r *Response) Json(data any) *Response {
	r.SetContentType("application/json")
	r.WriteHeader(r.statusCode)

	if err := json.NewEncoder(r).Encode(data); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
	}

	return r
}

func (r *Response) Text(text string) (int, error) {
	r.SetContentType("text/plain; charset=utf-8")
	r.WriteHeader(r.statusCode)
	return r.Write([]byte(text))
}
