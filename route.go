package router

import (
	"net/http"
	"net/url"
	"regexp"
)

type paramsList []string

type paramsValues map[string]string

type paramsRequirements map[string]*regexp.Regexp

type route struct {
	name               string
	action             Action
	priority           int
	method             string
	domain             string
	forwardRegexp      *regexp.Regexp
	reversePath        string
	requiredParams     paramsList
	paramsRequirements paramsRequirements
	defaultParams      paramsValues
}

var _ Route = &route{}

func (r *route) Priority() int {
	return r.priority
}

func (r *route) FindRoute(request *http.Request) (Route, bool) {
	if request.URL == nil {
		return nil, false
	}

	if !r.matchesDomain(request.URL) {
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

func (r *route) matchesDomain(requestURL *url.URL) bool {
	return r.domain == "" || r.domain == requestURL.Host
}

func (r *route) matchesMethod(request *http.Request) bool {
	return r.method == "" || r.method == request.Method
}

func (r *route) matchesPath(requestURL *url.URL) bool {
	return len(r.forwardRegexp.FindAllStringSubmatch(requestURL.Path, 1)) == 1
}
