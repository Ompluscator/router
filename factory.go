package router

import (
	"fmt"
	"regexp"
	"strings"
)

type factory struct {
	paramMatcher *regexp.Regexp
	requirement  *regexp.Regexp
}

func newFactory(paramMatcher *regexp.Regexp, requirement *regexp.Regexp) *factory {
	return &factory{
		paramMatcher: paramMatcher,
		requirement:  requirement,
	}
}

func (f *factory) createRoute(name string, path string, action Action, options RouteOptions) (*route, error) {
	required := f.createRequiredParams(path)
	defaults := f.createDefaultParams(options.DefaultParams)

	requirements, err := f.createParamsRequirements(name, required, options.ParamsRequirements)
	if err != nil {
		return nil, err
	}

	forward, err := f.createForwardRouteRegexp(name, path, required, requirements)
	if err != nil {
		return nil, err
	}

	return &route{
		name:               name,
		action:             name,
		priority:           options.Priority,
		method:             options.Method,
		secure:             options.Secure,
		host:               options.Host,
		forwardRegexp:      forward,
		reversePath:        path,
		requiredParams:     required,
		paramsRequirements: requirements,
		defaultParams:      defaults,
	}, nil
}

func (f *factory) createRouteGroup(name string, path string, options RouteGroupOptions) (*routeGroup, error) {
	required := f.createRequiredParams(path)
	defaults := f.createDefaultParams(options.DefaultParams)

	requirements, err := f.createParamsRequirements(name, required, options.ParamsRequirements)
	if err != nil {
		return nil, err
	}

	forward, err := f.createForwardRouteGroupRegexp(name, path, required)
	if err != nil {
		return nil, err
	}

	return &routeGroup{
		name:               name,
		secure:             options.Secure,
		host:               options.Host,
		forwardRegexp:      forward,
		reversePath:        path,
		paramsRequirements: requirements,
		defaultParams:      defaults,
		routes:             []routeFinder{},
		factory:            f,
	}, nil
}

func (f *factory) createParamsRequirements(name string, required paramsList, requirements map[string]string) (paramsRequirements, error) {
	result := paramsRequirements{}

	for _, key := range required {
		result[key] = f.requirement

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

func (f *factory) createForwardRouteRegexp(name string, path string, required paramsList, requirements paramsRequirements) (*regexp.Regexp, error) {
	forward := path

	for _, key := range required {
		compiled, ok := requirements[key]
		if !ok {
			compiled = f.requirement
		}

		wrapped := fmt.Sprintf("{%s}", key)
		forward = strings.Replace(forward, wrapped, compiled.String(), 1)
	}

	forward = fmt.Sprintf("^%s$", forward)
	result, err := regexp.Compile(forward)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for path "%s" in route "%s": %w`, path, name, err)
	}

	return result, nil
}

func (f *factory) createForwardRouteGroupRegexp(name string, path string, required paramsList) (*regexp.Regexp, error) {
	forward := path
	defaults := f.requirement.String()

	for _, key := range required {
		wrapped := fmt.Sprintf("{%s}", key)
		forward = strings.Replace(forward, wrapped, defaults, 1)
	}

	forward = fmt.Sprintf("^%s", forward)
	result, err := regexp.Compile(forward)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for path "%s" in route group "%s": %w`, path, name, err)
	}

	return result, nil
}

func (f *factory) createDefaultParams(defaults map[string]string) paramsValues {
	if defaults == nil {
		return paramsValues{}
	}
	return defaults
}

func (f *factory) createRequiredParams(path string) paramsList {
	return f.paramMatcher.FindAllString(path, -1)
}
