package router

import (
	"net/http"
	"net/url"
)

const (
	DefaultParamRequirement = `([^\/]+)`
)

type Action interface{}

type Options struct {
	Priority      int
	Secure        bool
	Host          string
	DefaultParams ParamsMap
}

type Route interface {
	Priority() int
	Name() string
	Path() string
	Action() Action
	URL(params ParamsMap) (*url.URL, error)
	ExtractParams(request *http.Request) (ParamsMap, error)
}

type RouteGroup interface {
	AddRoute(name string, path string, method string, action Action, options Options) error
	AddDeleteRoute(name string, path string, action Action) error
	AddGetRoute(name string, path string, action Action) error
	AddHeadRoute(name string, path string, action Action) error
	AddOptionsRoute(name string, path string, action Action) error
	AddPatchRoute(name string, path string, action Action) error
	AddPostRoute(name string, path string, action Action) error
	AddPutRoute(name string, path string, action Action) error
	AddRouteGroup(name string, path string, options Options) (RouteGroup, error)
}

type Router interface {
	RouteGroup
	FindRouteByRequest(request *http.Request) (Route, bool)
	FindRouteByName(name string) (Route, bool)
}

type Builder interface {
	SetSecure(secure bool) Builder
	SetHost(host string) Builder
	SetDefaultParamRequirement(expr string) Builder
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
