// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"errors"
)

var (
	NoHostError                = errors.New("no host")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

// Factory .
type Factory func([]string) Balancer

type Algorithm string

var factories = make(map[Algorithm]Factory)

func Build(algo Algorithm, hosts []string) (Balancer, error) {
	factory, ok := factories[algo]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}
	return factory(hosts), nil
}
