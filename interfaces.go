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

type ParamsMap map[string]string

func (p ParamsMap) Extend(other ParamsMap) ParamsMap {
	result := ParamsMap{}

	if p != nil {
		for k, v := range p {
			result[k] = v
		}
	}

	if other != nil {
		for k, v := range other {
			result[k] = v
		}
	}

	return result
}

type RouteOptions struct {
	Priority           int
	Method             string
	Secure             bool
	Host               string
	ParamsRequirements ParamsMap
	DefaultParams      ParamsMap
}

type RouteGroupOptions struct {
	Secure             bool
	Host               string
	ParamsRequirements ParamsMap
	DefaultParams      ParamsMap
}

type Route interface {
	Priority() int
	Name() string
	Path() string
	URL(params ParamsMap) (*url.URL, error)
	ExtractParams(request *http.Request) (ParamsMap, error)
}

type RouteGroup interface {
	AddRoute(name string, path string, action Action, options RouteOptions) error
	AddRouteGroup(name string, path string, options RouteGroupOptions) (RouteGroup, error)
}

type Router interface {
	RouteGroup
	FindRouteByRequest(request *http.Request) (Route, bool)
	FindRouteByName(name string) (Route, bool)
}

type Builder interface {
	SetSecure(secure bool) Builder
	SetHost(host string) Builder
	SetParamMatcher(expr string) Builder
	SetParamRequirement(expr string) Builder
	Build() (Router, error)
}

func New() Router {
	router, err := NewBuilder().Build()
	if err != nil {
		panic(err)
	}

	return router
}

func NewWithHost(secure bool, host string) Router {
	router, err := NewBuilder().SetSecure(secure).SetHost(host).Build()
	if err != nil {
		panic(err)
	}

	return router
}
