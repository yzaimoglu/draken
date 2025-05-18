package draken

import "net/http"

type Request struct {
	*http.Request
}

func GetRequest(r *http.Request) *Request {
	return &Request{r}
}

func (r *Request) RequestId() string {
	ctxReqId := r.Context().Value(ContextKeyRequestId)
	id, ok := ctxReqId.(string)
	if !ok {
		return ""
	}
	return id
}
