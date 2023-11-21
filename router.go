package suda

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var allMethods = []string{
	http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
	http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
	http.MethodTrace,
}

type Router struct {
	routes map[string][]ForwardTarget
}

func NewRouter() *Router {
	r := &Router{}
	r.routes = map[string][]ForwardTarget{}
	return r
}

func (r *Router) Add(uri string, target ForwardTarget) *Router {
	if r.routes[uri] == nil {
		r.routes[uri] = []ForwardTarget{}
	}
	r.routes[uri] = append(r.routes[uri], target)
	return r
}

func (r Router) Build(auth *AuthConfig) (handler http.Handler, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%s", r))
		}
	}()
	router := httprouter.New()
	for path, targets := range r.routes {
		forwarder := &Forwarder{
			Auth:    auth,
			Targets: targets,
		}
		for _, method := range allMethods {
			router.Handle(method, path, forwarder.serve)
		}
	}
	handler = router
	return
}
