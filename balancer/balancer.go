package balancer

import (
	"errors"
)

type Factory func([]string) Balancer

type Algorithm string

type Balancer interface {
	Add(string)
	Remove(string) bool
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

var factories = make(map[Algorithm]Factory)

func Build(algo Algorithm, hosts []string) (Balancer, error) {
	factory, ok := factories[algo]
	if !ok {
		return nil, errors.New("algorithm not supported")
	}
	return factory(hosts), nil
}
