package router

import (
	"fmt"
	"regexp"
)

type builder struct {
	secure           bool
	host             string
	paramMatcher     string
	paramRequirement string
}

func NewBuilder() Builder {
	return &builder{
		paramMatcher:     DefaultParamMatcher,
		paramRequirement: DefaultParamRequirement,
	}
}

func (b *builder) SetSecure(secure bool) Builder {
	b.secure = secure
	return b
}

func (b *builder) SetHost(host string) Builder {
	b.host = host
	return b
}

func (b *builder) SetParamMatcher(expr string) Builder {
	b.paramMatcher = expr
	return b
}

func (b *builder) SetParamRequirement(expr string) Builder {
	b.paramRequirement = expr
	return b
}

func (b *builder) Build() (Router, error) {
	paramMatcherCompiled, err := regexp.Compile(b.paramMatcher)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for param matcher "%s": %w`, b.paramMatcher, err)
	}

	paramRequirementCompiled, err := regexp.Compile(b.paramRequirement)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for param requirement "%s": %w`, b.paramRequirement, err)
	}

	factory := newFactory(paramMatcherCompiled, paramRequirementCompiled)

	group, err := factory.createRouteGroup("", "/", RouteGroupOptions{
		Secure: b.secure,
		Host:   b.host,
	})
	if err != nil {
		return nil, fmt.Errorf(`error while creating root group: %w`, err)
	}

	return &router{
		secure:  b.secure,
		host:    b.host,
		factory: factory,
		group:   group,
	}, nil
}
