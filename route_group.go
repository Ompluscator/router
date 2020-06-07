package router

import (
	"fmt"
	"net/http"
	"net/url"
	pathLib "path"
	"regexp"
	"strings"
)

type routeFinder interface {
	findRouteByRequest(request *http.Request) (Route, bool)
	findRouteByName(name string) (Route, bool)
}

type routeGroup struct {
	name               string
	secure             bool
	host               string
	forwardRegexp      *regexp.Regexp
	reversePath        string
	paramsRequirements paramsRequirements
	defaultParams      paramsValues
	routes             []routeFinder
	factory            *factory
}

var _ RouteGroup = &routeGroup{}
var _ routeFinder = &routeGroup{}

func (g *routeGroup) Name() string {
	return g.name
}

func (g *routeGroup) AddRoute(name string, path string, action Action, options RouteOptions) error {
	finalName := fmt.Sprintf("%s%s", g.getPrefix(), name)

	if _, ok := g.findRouteByName(finalName); ok {
		return fmt.Errorf(`route with name "%s" already exists`, finalName)
	}

	finalPath := pathLib.Join(g.reversePath, path)

	options.ParamsRequirements = g.paramsRequirements.toParamsMap().Extend(options.ParamsRequirements)
	options.DefaultParams = g.defaultParams.toParamsMap().Extend(options.DefaultParams)
	if g.secure {
		options.Secure = true
	}
	if g.host != "" && options.Host == "" {
		options.Host = g.host
	}

	route, err := g.factory.createRoute(finalName, finalPath, action, options)
	if err != nil {
		return err
	}

	g.routes = append(g.routes, route)

	return nil
}

func (g *routeGroup) AddRouteGroup(name string, path string, options RouteGroupOptions) (RouteGroup, error) {
	finalName := fmt.Sprintf("%s%s", g.getPrefix(), name)

	if _, ok := g.findRouteByName(finalName); ok {
		return nil, fmt.Errorf(`route with name "%s" already exists`, finalName)
	}

	finalPath := pathLib.Join(g.reversePath, path)

	options.ParamsRequirements = g.paramsRequirements.toParamsMap().Extend(options.ParamsRequirements)
	options.DefaultParams = g.defaultParams.toParamsMap().Extend(options.DefaultParams)
	if g.secure {
		options.Secure = true
	}
	if g.host != "" && options.Host == "" {
		options.Host = g.host
	}

	group, err := g.factory.createRouteGroup(finalName, finalPath, options)
	if err != nil {
		return nil, err
	}

	g.routes = append(g.routes, group)

	return group, nil
}

func (g *routeGroup) findRouteByRequest(request *http.Request) (Route, bool) {
	if request.URL == nil {
		return nil, false
	}

	if !g.matchesPath(request.URL) {
		return nil, false
	}

	var result Route

	for _, r := range g.routes {
		route, ok := r.findRouteByRequest(request)
		if !ok {
			continue
		}

		if result == nil || result.Priority() < route.Priority() {
			result = route
		}
	}

	return result, result != nil
}

func (g *routeGroup) findRouteByName(name string) (Route, bool) {
	if name == "" {
		return nil, false
	}

	prefix := g.getPrefix()

	if !strings.HasPrefix(name, prefix) {
		return nil, false
	}

	for _, r := range g.routes {
		route, ok := r.findRouteByName(name)
		if !ok {
			continue
		}

		return route, true
	}

	return nil, false
}

func (g *routeGroup) getPrefix() string {
	if g.name == "" {
		return ""
	}

	return fmt.Sprintf("%s.", g.name)
}

func (g *routeGroup) matchesPath(requestURL *url.URL) bool {
	return len(g.forwardRegexp.FindAllStringSubmatch(requestURL.Path, 1)) == 1
}
