package idpgo

import "net/http"

type HandlerWrapper func(w *ResponseWrapper, r *RequestWrapper)

func (h HandlerWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(&ResponseWrapper{w}, &RequestWrapper{r})
}