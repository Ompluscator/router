package router

import (
	"net/http"
	"net/url"
)

const (
	DefaultParamMatcher     = `\{([a-z]+)\}`
	DefaultParamRequirement = `(.+)`
)

type Action interface{}

type RouteOptions struct {
	Priority           int
	Method             string
	Secure             bool
	Host               string
	ParamsRequirements map[string]string
	DefaultParams      map[string]string
}

type Route interface {
	Priority() int
	Name() string
}

type RouteGroup interface {
	AddRoute(name string, path string, action Action, options RouteOptions) error
	AddRouteGroup(name string, path string) (RouteGroup, error)
}

type Router interface {
	RouteGroup
	FindRoute(request *http.Request) (Route, bool)
	RelativeURL(name string, values map[string]string) (*url.URL, error)
	AbsoluteURL(name string, values map[string]string) (*url.URL, error)
}
