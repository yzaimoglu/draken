package draken

import "net/http"

type Request struct {
	*http.Request
}

func GetRequest(r *http.Request) *Request {
	return &Request{r}
}

func (r *Request) CtxGetString(key ContextKey) string {
	ctxValReq := r.Context().Value(key)
	ctxVal, ok := ctxValReq.(string)
	if !ok {
		return ""
	}
	return ctxVal
}

func (r *Request) RequestId() string {
	return r.CtxGetString(ContextKeyRequestId)
}
