package router

import (
	"net/url"
)

type Action interface{}

type RouteOptions struct {
	Priority           int
	Method             string
	Domain             string
	ParamsRequirements map[string]string
	DefaultParams      map[string]string
}

type URLMatcher interface {
	matchesURLAddress(address url.URL) bool
	startsWithURLAddress(address url.URL) bool
}

type Route interface {
	URLMatcher
}

type RouteGroup interface {
	URLMatcher
	AddRoute(name string, path string, action Action, options RouteOptions) error
	AddRouteGroup(name string) (RouteGroup, error)
}

type Router interface {
	RouteGroup
	FindRouteForURLAddress(address url.URL) (Route, error)
	FindRouteForStringAddress(address url.URL) (Route, error)
	RelativeURL(name string, values map[string]string) (*url.URL, error)
	AbsoluteURL(name string, values map[string]string) (*url.URL, error)
}