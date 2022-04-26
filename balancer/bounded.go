// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "github.com/lafikl/consistent"

func init() {
	factories["bounded"] = NewBounded
}

// Bounded refers to consistent hash with bounded
type Bounded struct {
	ch *consistent.Consistent
}

// NewBounded create new Bounded balancer
func NewBounded(hosts []string) Balancer {
	c := &Bounded{consistent.New()}
	for _, h := range hosts {
		c.ch.Add(h)
	}
	return c
}

// Add new host to the balancer
func (c *Bounded) Add(host string) {
	c.ch.Add(host)
}

// Remove new host from the balancer
func (c *Bounded) Remove(host string) {
	c.ch.Remove(host)
}

// Balance selects a suitable host according to the key value
func (c *Bounded) Balance(key string) (string, error) {
	if len(c.ch.Hosts()) == 0 {
		return "", NoHostError
	}

	return c.ch.GetLeast(key)
}

// Inc refers to the number of connections to the server `+1`
func (c *Bounded) Inc(host string) {
	c.ch.Inc(host)
}

// Done refers to the number of connections to the server `-1`
func (c *Bounded) Done(host string) {
	c.ch.Done(host)
}
