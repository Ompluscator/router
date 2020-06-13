package router

import (
	"regexp"
)

type ParamsMap map[string]string

func (p ParamsMap) Extend(other ParamsMap) ParamsMap {
	result := ParamsMap{}

	if p != nil {
		for k, v := range p {
			result[k] = v
		}
	}

	if other != nil {
		for k, v := range other {
			result[k] = v
		}
	}

	return result
}

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
