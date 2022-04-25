// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "github.com/lafikl/consistent"

func init() {
	factories["bounded"] = NewBounded
}

type Bounded struct {
	ch *consistent.Consistent
}

func NewBounded(hosts []string) Balancer {
	c := &Bounded{consistent.New()}
	for _, h := range hosts {
		c.ch.Add(h)
	}
	return c
}

func (c *Bounded) Add(host string) {
	c.ch.Add(host)
}

func (c *Bounded) Remove(host string) {
	c.ch.Remove(host)
}

func (c *Bounded) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}

	return c.ch.GetLeast(key)
}

func (c *Bounded) Inc(host string) {
	c.ch.Inc(host)
}

func (c *Bounded) Done(host string) {
	c.ch.Done(host)
}
