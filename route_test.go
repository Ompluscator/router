package router

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"testing"
)

func TestRoute_Name(t *testing.T) {
	r := &route{}
	if r.Name() != "" {
		t.Errorf(`expected empty name but got "%s"`, r.Name())
	}

	r = &route{
		name: "name",
	}

	if r.Name() != "name" {
		t.Errorf(`expected name "name" but got "%s"`, r.Name())
	}
}

func TestRoute_Path(t *testing.T) {
	r := &route{}
	if r.Path() != "" {
		t.Errorf(`expected empty path but got "%s"`, r.Path())
	}

	r = &route{
		reversePath: "path",
	}

	if r.Path() != "path" {
		t.Errorf(`expected name "path" but got "%s"`, r.Path())
	}
}

func TestRoute_Priority(t *testing.T) {
	r := &route{}
	if r.Priority() != 0 {
		t.Errorf(`expected 0 priority but got %d`, r.Priority())
	}

	r = &route{
		priority: 100,
	}

	if r.Priority() != 100 {
		t.Errorf(`expected priority 100 but got %d`, r.Priority())
	}
}

func TestRoute_Action(t *testing.T) {
	r := &route{}
	if r.Action() != nil {
		t.Errorf(`expected nil but got %v`, r.Action())
	}

	action := func() {}
	r = &route{
		action: action,
	}
	if reflect.TypeOf(r.Action()) != reflect.TypeOf(action) {
		t.Errorf(`expected function but got %v`, reflect.TypeOf(r.Action()))
	}
}

func TestRoute_URL(t *testing.T) {
	cases := []struct {
		requiredParams paramsList
		requirements   paramsRequirements
		reversePath    string
		params         ParamsMap
		host           string
		secure         bool
		result         string
		missing        string
		invalid        string
	}{
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			reversePath:    "/{param1}/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
			host:    "",
			secure:  false,
			result:  "http:///value1/value2/value3",
			missing: "",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			reversePath:    "/{param1}/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
			host:    "domain.com",
			secure:  true,
			result:  "https://domain.com/value1/value2/value3",
			missing: "",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			reversePath:    "/{param1}/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
			},
			host:    "domain.com",
			secure:  true,
			result:  "",
			missing: "param3",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			reversePath: "/{param1}/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "de",
			},
			host:    "domain.com",
			secure:  true,
			result:  "https://domain.com/value/102/de",
			missing: "",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			reversePath: "/{param1}/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "gb",
			},
			host:    "domain.com",
			secure:  true,
			result:  "",
			missing: "",
			invalid: "param3",
		},
	}

	for _, c := range cases {
		r := route{
			host:               c.host,
			secure:             c.secure,
			requirement:        regexp.MustCompilePOSIX(DefaultParamRequirement),
			requiredParams:     c.requiredParams,
			reversePath:        c.reversePath,
			paramsRequirements: c.requirements,
		}

		result, err := r.URL(c.params)
		if c.missing == "" && c.invalid == "" {
			if err != nil {
				t.Errorf(`expected not to have error got: %v`, err)
			}
			if result.String() != c.result {
				t.Errorf(`expected result url "%s" but got "%s"`, c.result, result.String())
			}
		} else if c.missing != "" && (err == nil || err.Error() != fmt.Sprintf(`param "%s" is not provided`, c.missing)) {
			t.Errorf(`expected error for missing param "%s" but got %v`, c.missing, err)
		} else if c.invalid != "" && (err == nil || err.Error() != fmt.Sprintf(`invalid format provided for param "%s"`, c.invalid)) {
			t.Errorf(`expected error for missing param "%s" but got %v`, c.invalid, err)
		}
	}
}

func TestRoute_buildPath(t *testing.T) {
	cases := []struct {
		reversePath string
		params      ParamsMap
		result      string
	}{
		{
			reversePath: "/",
			params: ParamsMap{
				"param1": "value1",
			},
			result: "/",
		},
		{
			reversePath: "/{param1}",
			params: ParamsMap{
				"param1": "value1",
			},
			result: "/value1",
		},
		{
			reversePath: "/{param1}/{param2}",
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
			},
			result: "/value1/value2",
		},
		{
			reversePath: "/path/to/resource/{param1}/sub/resource/{param2}",
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
			},
			result: "/path/to/resource/value1/sub/resource/value2",
		},
		{
			reversePath: "/path/to/resource/{param1}/sub/resource/{param2}/{param3}",
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "de",
			},
			result: "/path/to/resource/value/sub/resource/102/de",
		},
	}

	for _, c := range cases {
		r := route{
			reversePath: c.reversePath,
		}

		result, err := r.buildPath(c.params)
		if err != nil {
			t.Errorf(`expected not to have error got: %v`, err)
		}

		if result != c.result {
			t.Errorf(`expected result path "%s" but got "%s"`, c.result, result)
		}
	}
}

func TestRoute_checkParams(t *testing.T) {
	cases := []struct {
		requiredParams paramsList
		requirements   paramsRequirements
		params         ParamsMap
		missing        string
		invalid        string
	}{
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
				"param3": "value2",
			},
			missing: "",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			params: ParamsMap{
				"param2": "value2",
				"param3": "value2",
			},
			missing: "param1",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			params: ParamsMap{
				"param1": "value1",
				"param3": "value2",
			},
			missing: "param2",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			params: ParamsMap{
				"param1": "value1",
				"param2": "value2",
			},
			missing: "param3",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			params:         ParamsMap{},
			missing:        "param1",
			invalid:        "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements:   paramsRequirements{},
			missing:        "param1",
			invalid:        "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "de",
			},
			missing: "",
			invalid: "",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			params: ParamsMap{
				"param1": "value1",
				"param2": "102",
				"param3": "de",
			},
			missing: "",
			invalid: "param1",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			params: ParamsMap{
				"param1": "value",
				"param2": "102a",
				"param3": "de",
			},
			missing: "",
			invalid: "param2",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
			},
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "gb",
			},
			missing: "",
			invalid: "param3",
		},
		{
			requiredParams: paramsList{"param1", "param2", "param3"},
			requirements: paramsRequirements{
				"param1": regexp.MustCompile(`([a-z]+)`),
				"param2": regexp.MustCompile(`([0-9]+)`),
				"param3": regexp.MustCompile(`(de|fr)`),
				"param4": regexp.MustCompile(`([a-z]+)`),
			},
			params: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "de",
				"param4": "102",
			},
			missing: "",
			invalid: "param4",
		},
	}

	for _, c := range cases {
		r := route{
			requirement:        regexp.MustCompilePOSIX(DefaultParamRequirement),
			requiredParams:     c.requiredParams,
			paramsRequirements: c.requirements,
		}

		err := r.checkParams(c.params)
		if c.missing == "" && c.invalid == "" && err != nil {
			t.Errorf(`expected not to have error got: %v`, err)
		} else if c.missing != "" && (err == nil || err.Error() != fmt.Sprintf(`param "%s" is not provided`, c.missing)) {
			t.Errorf(`expected error for missing param "%s" but got %v`, c.missing, err)
		} else if c.invalid != "" && (err == nil || err.Error() != fmt.Sprintf(`invalid format provided for param "%s"`, c.invalid)) {
			t.Errorf(`expected error for missing param "%s" but got %v`, c.invalid, err)
		}
	}
}

func TestRoute_ExtractParams(t *testing.T) {
	cases := []struct {
		forwardRegexp  *regexp.Regexp
		requiredParams paramsList
		url            string
		result         ParamsMap
		err            string
	}{
		{
			forwardRegexp:  regexp.MustCompile(`^/$`),
			requiredParams: paramsList{},
			url:            "/",
			result:         ParamsMap{},
			err:            "",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)$`),
			requiredParams: paramsList{"param1"},
			url:            "/path/to/value1",
			result: ParamsMap{
				"param1": "value1",
			},
			err: "",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)/([0-9]+)/(de|fr)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			url:            "/path/to/value/102/de",
			result: ParamsMap{
				"param1": "value",
				"param2": "102",
				"param3": "de",
			},
			err: "",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			url:            "/path/to/value1/value2/value3",
			result: ParamsMap{
				"param1": "value1",
				"param2": "value2",
				"param3": "value3",
			},
			err: "",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2"},
			url:            "/path/to/value1/value2/value3",
			result:         ParamsMap{},
			err:            "url does not belong to route",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			url:            "/path/to/value1/value2",
			result:         ParamsMap{},
			err:            "url does not belong to route",
		},
		{
			forwardRegexp:  regexp.MustCompile(`^/path/to/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			url:            "",
			result:         ParamsMap{},
			err:            "url is not provided",
		},
	}

	for _, c := range cases {
		r := &route{
			forwardRegexp:  c.forwardRegexp,
			requiredParams: c.requiredParams,
		}

		var err error
		var req *http.Request

		if c.url != "" {
			req, err = http.NewRequest(http.MethodGet, c.url, nil)
			if err != nil {
				t.Errorf(`not expected error but get %s`, err.Error())
			}
		}

		result, err := r.ExtractParams(req)
		if c.err != "" {
			if err == nil {
				t.Error("expected error but got nil")
			} else if err.Error() != c.err {
				t.Errorf(`expected error "%s" but got "%s"`, c.err, err.Error())
			}
		} else if !reflect.DeepEqual(c.result, result) {
			t.Errorf(`expected params %v but got %v`, c.result, result)
		}
	}
}

func TestRoute_findRouteByRequest(t *testing.T) {
	cases := []struct {
		requestMethod  string
		requestUrl     string
		secure         bool
		host           string
		method         string
		forwardRegexp  *regexp.Regexp
		requiredParams paramsList
		result         bool
	}{
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.de",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         true,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         true,
			host:           "domain.de",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         true,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.de",
			method:         "",
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         true,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         true,
		},
		{
			requestMethod:  http.MethodPost,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.de",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         false,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.com",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2", "param3"},
			result:         false,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.de",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2"},
			result:         false,
		},
		{
			requestMethod:  http.MethodGet,
			requestUrl:     "http://domain.de/value1/value2/value3",
			secure:         false,
			host:           "domain.de",
			method:         http.MethodGet,
			forwardRegexp:  regexp.MustCompile(`^/([^\/]+)/([^\/]+)$`),
			requiredParams: paramsList{"param1", "param2"},
			result:         false,
		},
	}

	for _, c := range cases {
		r := &route{
			secure:         c.secure,
			host:           c.host,
			method:         c.method,
			forwardRegexp:  c.forwardRegexp,
			requiredParams: c.requiredParams,
		}

		var err error
		var req *http.Request

		if c.requestUrl != "" {
			req, err = http.NewRequest(c.requestMethod, c.requestUrl, nil)
			if err != nil {
				t.Errorf(`not expected error but get %s`, err.Error())
			}
		}

		result, ok := r.findRouteByRequest(req)
		if ok != c.result {
			t.Errorf(`expected %t but got %t`, c.result, ok)
		} else if ok && !reflect.DeepEqual(result, r) {
			t.Errorf(`expected %v but got %v`, r, result)
		}
	}
}

func TestRoute_findRouteByName(t *testing.T) {
	cases := []struct {
		name      string
		routeName string
		result    bool
	}{
		{
			name:      "",
			routeName: "route",
			result:    false,
		},
		{
			name:      "wrong",
			routeName: "route",
			result:    false,
		},
		{
			name:      "route",
			routeName: "",
			result:    false,
		},
		{
			name:      "route",
			routeName: "route",
			result:    true,
		},
	}

	for _, c := range cases {
		r := &route{
			name: c.routeName,
		}

		result, ok := r.findRouteByName(c.name)
		if ok != c.result {
			t.Errorf(`expected %t but got %t`, c.result, ok)
		} else if ok && !reflect.DeepEqual(result, r) {
			t.Errorf(`expected %v but got %v`, r, result)
		}
	}
}

func TestRoute_matchesHost(t *testing.T) {
	cases := []struct {
		host    string
		urlHost string
		result  bool
	}{
		{
			host:    "",
			urlHost: "",
			result:  true,
		},
		{
			host:    "",
			urlHost: "host",
			result:  true,
		},
		{
			host:    "host",
			urlHost: "",
			result:  false,
		},
		{
			host:    "host",
			urlHost: "wrong",
			result:  false,
		},
		{
			host:    "host",
			urlHost: "host",
			result:  true,
		},
	}

	for _, c := range cases {
		r := route{
			host: c.host,
		}

		u := &url.URL{
			Host: c.urlHost,
		}

		result := r.matchesHost(u)
		if result != c.result {
			t.Errorf(`expecte %t but got %t`, c.result, result)
		}
	}
}

func TestRoute_matchesMethod(t *testing.T) {
	cases := []struct {
		method        string
		requestMethod string
		result        bool
	}{
		{
			method:        "",
			requestMethod: "",
			result:        true,
		},
		{
			method:        "",
			requestMethod: http.MethodGet,
			result:        true,
		},
		{
			method:        http.MethodGet,
			requestMethod: "",
			result:        false,
		},
		{
			method:        http.MethodGet,
			requestMethod: http.MethodPost,
			result:        false,
		},
		{
			method:        http.MethodGet,
			requestMethod: http.MethodGet,
			result:        true,
		},
	}

	for _, c := range cases {
		r := route{
			method: c.method,
		}

		req := &http.Request{
			Method: c.requestMethod,
		}

		result := r.matchesMethod(req)
		if result != c.result {
			t.Errorf(`expecte %t but got %t`, c.result, result)
		}
	}
}

func TestRoute_matchesPath(t *testing.T) {
	cases := []struct {
		forwardRegexp *regexp.Regexp
		requiredParam paramsList
		urlPath       string
		result        bool
	}{
		{
			forwardRegexp: regexp.MustCompile(`^/$`),
			requiredParam: paramsList{},
			urlPath:       "/",
			result:        true,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/$`),
			requiredParam: paramsList{"param1"},
			urlPath:       "/",
			result:        false,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([^\/]+)$`),
			requiredParam: paramsList{"param1"},
			urlPath:       "/value1",
			result:        true,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)/([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2", "param3"},
			urlPath:       "/value/102/de",
			result:        true,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)/([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2", "param3"},
			urlPath:       "/value/102/gb",
			result:        false,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)/([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2", "param3"},
			urlPath:       "/value/a/de",
			result:        false,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)/([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2", "param3"},
			urlPath:       "/value1/102/de",
			result:        false,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)/([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2"},
			urlPath:       "/value/102/de",
			result:        false,
		},
		{
			forwardRegexp: regexp.MustCompile(`^/([a-z]+)-([0-9]+)/(de|fr)$`),
			requiredParam: paramsList{"param1", "param2", "param3"},
			urlPath:       "/value-102/de",
			result:        true,
		},
	}

	for _, c := range cases {
		r := route{
			forwardRegexp:  c.forwardRegexp,
			requiredParams: c.requiredParam,
		}

		u := &url.URL{
			Path: c.urlPath,
		}

		result := r.matchesPath(u)
		if result != c.result {
			t.Errorf(`expecte %t but got %t`, c.result, result)
		}
	}
}
