package router

import (
	"net/http"
)

type router struct {
	factory *factory
	group   *routeGroup
}

var _ Router = &router{}

func (r *router) AddRoute(name string, path string, method string, action Action, options Options) error {
	return r.group.AddRoute(name, path, method, action, options)
}

func (r *router) AddDeleteRoute(name string, path string, action Action) error {
	return r.group.AddDeleteRoute(name, path, action)
}

func (r *router) AddGetRoute(name string, path string, action Action) error {
	return r.group.AddGetRoute(name, path, action)
}

func (r *router) AddHeadRoute(name string, path string, action Action) error {
	return r.group.AddHeadRoute(name, path, action)
}

func (r *router) AddOptionsRoute(name string, path string, action Action) error {
	return r.group.AddOptionsRoute(name, path, action)
}

func (r *router) AddPatchRoute(name string, path string, action Action) error {
	return r.group.AddPatchRoute(name, path, action)
}

func (r *router) AddPostRoute(name string, path string, action Action) error {
	return r.group.AddPostRoute(name, path, action)
}

func (r *router) AddPutRoute(name string, path string, action Action) error {
	return r.group.AddPutRoute(name, path, action)
}

func (r *router) AddRouteGroup(name string, path string, options Options) (RouteGroup, error) {
	return r.group.AddRouteGroup(name, path, options)
}

func (r *router) FindRouteByRequest(request *http.Request) (Route, bool) {
	return r.group.findRouteByRequest(request)
}

func (r *router) FindRouteByName(name string) (Route, bool) {
	return r.group.findRouteByName(name)
}