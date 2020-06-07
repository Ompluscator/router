package router

import (
	"net/http"
	"net/url"
)

type router struct {
	secure  bool
	host    string
	factory *factory
	group   *routeGroup
}

var _ Router = &router{}

func (r *router) AddRoute(name string, path string, action Action, options RouteOptions) error {
	return r.group.AddRoute(name, path, action, options)
}

func (r *router) AddRouteGroup(name string, path string, options RouteGroupOptions) (RouteGroup, error) {
	return r.group.AddRouteGroup(name, path, options)
}

func (r *router) FindRouteByRequest(request *http.Request) (Route, bool) {
	return r.group.findRouteByRequest(request)
}

func (r *router) FindRouteByName(name string) (Route, bool) {
	return r.group.findRouteByName(name)
}