// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"math/rand"
	"sync"
	"time"
)

func init() {
	factories["random"] = NewRandom
}

// Random will randomly select a http server from the server
type Random struct {
	sync.Mutex
	hosts []string
	rnd   *rand.Rand
}

// NewRandom create new Random balancer
func NewRandom(hosts []string) Balancer {
	return &Random{hosts: hosts,
		rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

// Add new host to the balancer
func (r *Random) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

// Remove new host from the balancer
func (r *Random) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

// Balance selects a suitable host according
func (r *Random) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}

// Inc .
func (r *Random) Inc(_ string) {}

// Done .
func (r *Random) Done(_ string) {}
