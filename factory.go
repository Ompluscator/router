package router

import (
	"fmt"
	"regexp"
	"strings"
)

const defaultParamMatcher = `\{([a-z]+[\:]{0,1}[^\}]*)\}`

type factory struct {
	paramMatcher *regexp.Regexp
	requirement  *regexp.Regexp
}

func newFactory(requirement *regexp.Regexp) *factory {
	paramMatcher := regexp.MustCompile(defaultParamMatcher)

	return &factory{
		paramMatcher: paramMatcher,
		requirement:  requirement,
	}
}

func (f *factory) createRoute(name string, path string, method string, action Action, options Options) (*route, error) {
	pairs, required, err := f.createParams(path)
	if err != nil {
		return nil, err
	}

	defaults := f.createDefaultParams(options.DefaultParams)

	requirements, err := f.createParamsRequirements(name, pairs)
	if err != nil {
		return nil, err
	}

	forward, err := f.createForwardRouteRegexp(name, path, required, requirements)
	if err != nil {
		return nil, err
	}

	return &route{
		name:               name,
		action:             action,
		priority:           options.Priority,
		method:             method,
		secure:             options.Secure,
		host:               options.Host,
		forwardRegexp:      forward,
		reversePath:        path,
		requiredParams:     required,
		paramsRequirements: requirements,
		defaultParams:      defaults,
		requirement:        f.requirement,
	}, nil
}

func (f *factory) createRouteGroup(name string, path string, options Options) (*routeGroup, error) {
	pairs, required, err := f.createParams(path)
	if err != nil {
		return nil, err
	}

	defaults := f.createDefaultParams(options.DefaultParams)

	requirements, err := f.createParamsRequirements(name, pairs)
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
		originalPath:       path,
		paramsRequirements: requirements,
		defaultParams:      defaults,
		routes:             []routeFinder{},
		factory:            f,
	}, nil
}

func (f *factory) createParamsRequirements(name string, requirements map[string]string) (paramsRequirements, error) {
	result := paramsRequirements{}

	for key, value := range requirements {
		compiled, err := regexp.Compile(value)
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

		wrapped := fmt.Sprintf("{%s}", requirements[key].String())
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

func (f *factory) createParams(path string) (ParamsMap, paramsList, error) {
	paramMap := ParamsMap{}
	var list paramsList

	params := f.paramMatcher.FindAllStringSubmatch(path, -1)

	for _, v := range params {
		if len(v) != 2 {
			return nil, nil, fmt.Errorf(`invlaid path provided: %s`, path)
		}

		matches := strings.Split(v[1], ":")
		if len(matches) == 0 || matches[0] == "" {
			return nil, nil, fmt.Errorf(`empty param is provided in path: %s`, path)
		}

		if _, ok := paramMap[matches[0]]; ok {
			return nil, nil, fmt.Errorf(`param with name "%s" is provided mutiple times in path: %s`, matches[0], path)
		}

		requirement := f.requirement.String()
		if len(matches) > 1 && matches[1] != "" {
			requirement = fmt.Sprintf("(%s)", matches[1])
		}

		paramMap[matches[0]] = requirement
		list = append(list, matches[0])
	}

	return paramMap, list, nil
}
