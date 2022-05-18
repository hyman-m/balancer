// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "github.com/lafikl/consistent"

func init() {
	factories[BoundedBalancer] = NewBounded
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
func (b *Bounded) Add(host string) {
	b.ch.Add(host)
}

// Remove new host from the balancer
func (b *Bounded) Remove(host string) {
	b.ch.Remove(host)
}

// Balance selects a suitable host according to the key value
func (b *Bounded) Balance(key string) (string, error) {
	if len(b.ch.Hosts()) == 0 {
		return "", NoHostError
	}
	return b.ch.GetLeast(key)
}

// Inc refers to the number of connections to the server `+1`
func (b *Bounded) Inc(host string) {
	b.ch.Inc(host)
}

// Done refers to the number of connections to the server `-1`
func (b *Bounded) Done(host string) {
	b.ch.Done(host)
}
