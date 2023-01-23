// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"math/rand"
	"time"
)

func init() {
	factories[RandomBalancer] = NewRandom
}

// Random will randomly select a http server from the server
type Random struct {
	BaseBalancer
	rnd *rand.Rand
}

// NewRandom create new Random balancer
func NewRandom(hosts []string) Balancer {
	return &Random{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
		rnd: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Balance selects a suitable host according
func (r *Random) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	return r.hosts[r.rnd.Intn(len(r.hosts))], nil
}
