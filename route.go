package router

import (
	"net/http"
	"net/url"
	"regexp"
)

type paramsList []string

type paramsValues map[string]string

func (p paramsValues) toParamsMap() ParamsMap {
	result := ParamsMap{}

	if p != nil {
		for k, v := range p {
			result[k] = v
		}
	}

	return result
}

type paramsRequirements map[string]*regexp.Regexp

func (p paramsRequirements) toParamsMap() ParamsMap {
	result := ParamsMap{}

	if p != nil {
		for k, v := range p {
			result[k] = v.String()
		}
	}

	return result
}

type route struct {
	name               string
	action             Action
	priority           int
	method             string
	secure             bool
	host               string
	forwardRegexp      *regexp.Regexp
	reversePath        string
	requiredParams     paramsList
	paramsRequirements paramsRequirements
	defaultParams      paramsValues
}

var _ Route = &route{}
var _ routeFinder = &route{}

func (r *route) Name() string {
	return r.name
}

func (r *route) Priority() int {
	return r.priority
}

func (r *route) findRouteByRequest(request *http.Request) (Route, bool) {
	if request.URL == nil {
		return nil, false
	}

	if !r.matchesHost(request.URL) {
		return nil, false
	}

	if !r.matchesMethod(request) {
		return nil, false
	}

	if !r.matchesPath(request.URL) {
		return nil, false
	}

	return r, true
}

func (r *route) findRouteByName(name string) (Route, bool) {
	if name == "" {
		return nil, false
	}

	if r.name != name {
		return nil, false
	}

	return r, true
}

func (r *route) matchesHost(requestURL *url.URL) bool {
	return r.host == "" || r.host == requestURL.Host
}

func (r *route) matchesMethod(request *http.Request) bool {
	return r.method == "" || r.method == request.Method
}

func (r *route) matchesPath(requestURL *url.URL) bool {
	return len(r.forwardRegexp.FindAllStringSubmatch(requestURL.Path, 1)) == 1
}
