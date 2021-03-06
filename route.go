package router

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

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
	requirement        *regexp.Regexp
}

var _ Route = &route{}
var _ routeFinder = &route{}

func (r *route) Name() string {
	return r.name
}

func (r *route) Priority() int {
	return r.priority
}

func (r *route) Path() string {
	return r.reversePath
}

func (r *route) Action() Action {
	return r.action
}

func (r *route) URL(params ParamsMap) (*url.URL, error) {
	finalParams := r.defaultParams.toParamsMap().Extend(params)

	err := r.checkParams(finalParams)
	if err != nil {
		return nil, err
	}

	path, err := r.buildPath(finalParams)
	if err != nil {
		return nil, err
	}

	scheme := "http"
	if r.secure {
		scheme = "https"
	}

	return &url.URL{
		Scheme: scheme,
		Host:   r.host,
		Path:   path,
	}, nil
}

func (r *route) ExtractParams(request *http.Request) (ParamsMap, error) {
	if request == nil || request.URL == nil {
		return nil, errors.New("url is not provided")
	}

	matches, err := r.getMatchesPath(request.URL)
	if err != nil {
		return nil, err
	}

	result := ParamsMap{}
	for index, key := range r.requiredParams {
		result[key] = matches[0][index+1]
	}

	return result, nil
}

func (r *route) buildPath(params ParamsMap) (string, error) {
	path := r.reversePath

	for key, value := range params {
		wrapped := fmt.Sprintf("{%s}", key)
		if strings.Index(path, wrapped) == -1 {
			continue
		}

		path = strings.Replace(path, wrapped, value, 1)
	}

	return path, nil
}

func (r *route) checkParams(params ParamsMap) error {
	checked := map[string]bool{}

	for _, key := range r.requiredParams {
		checked[key] = true

		if params == nil {
			return fmt.Errorf(`param "%s" is not provided`, key)
		}

		value, ok := params[key]
		if !ok {
			return fmt.Errorf(`param "%s" is not provided`, key)
		}

		requirement, ok := r.paramsRequirements[key]
		if !ok {
			requirement = r.requirement
		}

		matches := requirement.FindAllString(value, 1)
		if len(matches) == 0 || matches[0] != value {
			return fmt.Errorf(`invalid format provided for param "%s"`, key)
		}
	}

	for key, compiled := range r.paramsRequirements {
		if checked[key] {
			continue
		}

		value, ok := params[key]
		if !ok {
			continue
		}

		matches := compiled.FindAllString(value, 1)
		if len(matches) == 0 || matches[0] != value {
			return fmt.Errorf(`invalid format provided for param "%s"`, key)
		}
	}

	return nil
}

func (r *route) findRouteByRequest(request *http.Request) (Route, bool) {
	if request == nil || request.URL == nil {
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
	_, err := r.getMatchesPath(requestURL)
	return err == nil
}

func (r *route) getMatchesPath(requestURL *url.URL) ([][]string, error) {
	matches := r.forwardRegexp.FindAllStringSubmatch(requestURL.Path, 1)
	if len(matches) != 1 || len(matches[0]) != len(r.requiredParams)+1 {
		return nil, errors.New("url does not belong to route")
	}

	return matches, nil
}
