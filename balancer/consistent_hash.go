// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"github.com/lafikl/consistent"
)

func init() {
	factories["consistent-hash"] = NewConsistent
}

type Consistent struct {
	ch *consistent.Consistent
}

func NewConsistent(hosts []string) Balancer {
	c := &Consistent{consistent.New()}
	for _, h := range hosts {
		c.ch.Add(h)
	}
	return c
}

func (c *Consistent) Add(host string) {
	c.ch.Add(host)
}

func (c *Consistent) Remove(host string) {
	c.ch.Remove(host)
}

func (c *Consistent) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}

	return c.ch.Get(key)
}

func (c *Consistent) Inc(_ string) {}

func (c *Consistent) Done(_ string) {}
