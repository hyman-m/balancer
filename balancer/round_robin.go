// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import "sync/atomic"

//RoundRobin will select the server in turn from the server to proxy
type RoundRobin struct {
	BaseBalancer
	i atomic.Uint64
}

func init() {
	factories[R2Balancer] = NewRoundRobin
}

// NewRoundRobin create new RoundRobin balancer
func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		i: atomic.Uint64{},
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

// Balance selects a suitable host according
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	host := r.hosts[r.i.Add(1)%uint64(len(r.hosts))]
	return host, nil
}
