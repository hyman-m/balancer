// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sync"
)

//RoundRobin will select the server in turn from the server to proxy
type RoundRobin struct {
	sync.RWMutex
	i     uint64
	hosts []string
}

func init() {
	factories[R2Balancer] = NewRoundRobin
}

// NewRoundRobin create new RoundRobin balancer
func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{i: 0, hosts: hosts}
}

// Add new host to the balancer
func (r *RoundRobin) Add(host string) {
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
func (r *RoundRobin) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return
		}
	}
}

// Balance selects a suitable host according
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}

// Inc .
func (r *RoundRobin) Inc(_ string) {}

// Done .
func (r *RoundRobin) Done(_ string) {}
