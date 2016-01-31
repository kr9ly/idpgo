package idpgo

import (
	"net/http"
	"html/template"
)

type ResponseWrapper struct {
	http.ResponseWriter
}

func (w *ResponseWrapper) BadRequest(m string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(m))
}

func (w *ResponseWrapper) MethodNotAllowed(m string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte(m))
}

func (w *ResponseWrapper) ResponseTemplate(t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	t.Execute(w, data)
}

func (w *ResponseWrapper) SeeOther(uri string) {
	w.Header().Set("Location", uri)
	w.WriteHeader(http.StatusSeeOther)
}