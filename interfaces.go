package router

import (
	"net/url"
	"regexp"
)

type Action interface{}

type ParamsRequirements map[string]*regexp.Regexp

type ParamsValues map[string]string

type RouteOptions struct {
	Priority           int
	Method             string
	Domain             string
	ParamsRequirements ParamsRequirements
	DefaultParams      ParamsValues
}

type URLMatcher interface {
	MatchesURLAddress(address url.URL) bool
	StartsWithURLAddress(address url.URL) bool
	MatchStringAddress(address string) bool
	StartsWithStringAddress(address string) bool
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
	RelativeURL(name string, values ParamsValues) (*url.URL, error)
	AbsoluteURL(name string, values ParamsValues) (*url.URL, error)
}