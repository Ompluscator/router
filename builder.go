package router

import (
	"fmt"
	"regexp"
)

type builder struct {
	secure           bool
	host             string
	paramRequirement string
}

func NewBuilder() Builder {
	return &builder{
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

func (b *builder) SetDefaultParamRequirement(expr string) Builder {
	b.paramRequirement = expr
	return b
}

func (b *builder) Build() (Router, error) {
	paramRequirementCompiled, err := regexp.Compile(b.paramRequirement)
	if err != nil {
		return nil, fmt.Errorf(`error while compiling regexp for param requirement "%s": %w`, b.paramRequirement, err)
	}

	factory := newFactory(paramRequirementCompiled)

	group, err := factory.createRouteGroup("", "/", Options{
		Secure: b.secure,
		Host:   b.host,
	})
	if err != nil {
		return nil, fmt.Errorf(`error while creating root group: %w`, err)
	}

	return &router{
		factory: factory,
		group:   group,
	}, nil
}
