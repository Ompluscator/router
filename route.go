package router

import (
	"fmt"
	"regexp"
	"strings"
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

type routeFactory struct {
	paramRegexp   *regexp.Regexp
	defaultRegexp *regexp.Regexp
}

func newRouteFactory() (*routeFactory, error) {
	paramCompiled, err := regexp.Compile(`\{([a-z]+)\}`)
	if err != nil {
		return nil, err
	}

	defaultCompiled, err := regexp.Compile(`(.+)`)
	if err != nil {
		return nil, err
	}

	return &routeFactory{
		paramRegexp:   paramCompiled,
		defaultRegexp: defaultCompiled,
	}, nil
}

func (f *routeFactory) create(name string, path string, action Action, options RouteOptions) (*route, error) {
	required := f.createRequiredParams(path)
	defaults := f.createDefaultParams(options.DefaultParams)

	requirements, err := f.createParamsRequirements(name, required, options.ParamsRequirements)
	if err != nil {
		return nil, err
	}

	forward, err := f.createForwardRegexp(name, path, required, requirements)
	if err != nil {
		return nil, err
	}

	return &route{
		name:               name,
		action:             name,
		priority:           options.Priority,
		method:             options.Method,
		domain:             options.Domain,
		forwardRegexp:      forward,
		reversePath:        path,
		requiredParams:     required,
		paramsRequirements: requirements,
		defaultParams:      defaults,
	}, nil
}

func (f *routeFactory) createParamsRequirements(name string, required paramsList, requirements map[string]string) (paramsRequirements, error) {
	result := paramsRequirements{}

	for _, key := range required {
		result[key] = f.defaultRegexp

		if requirements == nil {
			continue
		}

		regexpString, ok := requirements[key]
		if !ok {
			continue
		}

		compiled, err := regexp.Compile(fmt.Sprintf("(%s)", regexpString))
		if err != nil {
			return nil, fmt.Errorf(`error while compiling regexp for param "%s" in route "%s": %w`, key, name, err)
		}

		result[key] = compiled
	}

	return result, nil
}

func (f *routeFactory) createForwardRegexp(name string, path string, required paramsList, requirements paramsRequirements) (*regexp.Regexp, error) {
	forward := path

	for _, key := range required {
		compiled, ok := requirements[key]
		if !ok {
			compiled = f.defaultRegexp
		}

		wrapped := fmt.Sprintf("{%s}", key)
		forward = strings.Replace(forward, wrapped, compiled.String(), 1)
	}

	result, err := regexp.Compile(forward)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for path "%s" in route "%s": %w`, path, name, err)
	}

	return result, nil
}

func (f *routeFactory) createDefaultParams(defaults map[string]string) paramsValues {
	if defaults == nil {
		return paramsValues{}
	}
	return defaults
}

func (f *routeFactory) createRequiredParams(path string) paramsList {
	return f.paramRegexp.FindAllString(path, -1)
}
