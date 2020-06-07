package router

import (
	"net/http"
	"net/url"
)

const (
	DefaultParamMather      = `\{([a-z]+)\}`
	DefaultParamRequirement = `(.+)`
)

type Action interface{}

type RouteOptions struct {
	Priority           int
	Method             string
	Domain             string
	ParamsRequirements map[string]string
	DefaultParams      map[string]string
}

type RouteFinder interface {
	FindRoute(request *http.Request) (Route, bool)
}

type Route interface {
	RouteFinder
	Priority() int
}

type RouteGroup interface {
	RouteFinder
	AddRoute(name string, path string, action Action, options RouteOptions) error
	AddRouteGroup(name string, path string) (RouteGroup, error)
}

type Router interface {
	RouteGroup
	RelativeURL(name string, values map[string]string) (*url.URL, error)
	AbsoluteURL(name string, values map[string]string) (*url.URL, error)
}
