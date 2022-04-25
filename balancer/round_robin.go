// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sync"
)

type RoundRobin struct {
	sync.Mutex
	i     uint64
	hosts []string
}

func init() {
	factories["round-robin"] = NewRoundRobin
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{i: 0, hosts: hosts}
}

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

func (r *RoundRobin) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
		}
	}
}

func (r *RoundRobin) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}

func (r *RoundRobin) Inc(_ string) {}

func (r *RoundRobin) Done(_ string) {}
