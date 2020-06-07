package router

import (
	"net/http"
	"net/url"
	pathLib "path"
	"regexp"
)

type routeGroup struct {
	name          string
	forwardRegexp *regexp.Regexp
	reversePath   string
	routes        []RouteFinder
	factory       *factory
}

var _ RouteGroup = &routeGroup{}

func (g *routeGroup) FindRoute(request *http.Request) (Route, bool) {
	if request.URL == nil {
		return nil, false
	}

	if !g.matchesPath(request.URL) {
		return nil, false
	}

	var result Route

	for _, r := range g.routes {
		route, ok := r.FindRoute(request)
		if !ok {
			continue
		}

		if result == nil || result.Priority() < route.Priority() {
			result = route
		}
	}

	return result, result != nil
}

func (g *routeGroup) AddRoute(name string, path string, action Action, options RouteOptions) error {
	final := pathLib.Join(g.reversePath, path)

	route, err := g.factory.createRoute(name, final, action, options)
	if err != nil {
		return err
	}

	g.routes = append(g.routes, route)

	return nil
}

func (g *routeGroup) AddRouteGroup(name string, path string) (RouteGroup, error) {
	final := pathLib.Join(g.reversePath, path)

	group, err := g.factory.createRouteGroup(name, final)
	if err != nil {
		return nil, err
	}

	g.routes = append(g.routes, group)

	return group, nil
}

func (g *routeGroup) matchesPath(requestURL *url.URL) bool {
	return len(g.forwardRegexp.FindAllStringSubmatch(requestURL.Path, 1)) == 1
}
