package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
	Validator StructValidator
}

func NewRouter() *Router {
	return &Router{
		Router: mux.NewRouter(),
	}
}

func NewRouterWithCustomValidator(validator StructValidator) *Router {
	return &Router{
		Router:    mux.NewRouter(),
		Validator: validator,
	}
}

func (r *Router) SetValidator(validator StructValidator) *Router {
	r.Validator = validator
	return r
}

func (r *Router) Group(path string) *Router {
	return &Router{
		Router:    r.PathPrefix(path).Subrouter(),
		Validator: r.Validator,
	}
}

func (r *Router) GET(path string, f HandleFunc) {
	r.Handle("GET", path, f)
}

func (r *Router) POST(path string, f HandleFunc) {
	r.Handle("POST", path, f)
}

func (r *Router) PUT(path string, f HandleFunc) {
	r.Handle("PUT", path, f)
}

func (r *Router) PATCH(path string, f HandleFunc) {
	r.Handle("PATCH", path, f)
}

func (r *Router) DELETE(path string, f HandleFunc) {
	r.Handle("DELETE", path, f)
}

func (r *Router) OPTIONS(path string, f HandleFunc) {
	r.Handle("OPTIONS", path, f)
}

func (r *Router) Handle(method string, path string, f HandleFunc) {
	r.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		f(&Context{
			Writer:    w,
			Request:   req,
			Validator: r.Validator,
		})
	}).Methods(method)
}
