// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"hash/crc32"
	"sync"
)

func init() {
	factories[IPHashBalancer] = NewIPHash
}

// IPHash will choose a host based on the client's IP address
type IPHash struct {
	sync.RWMutex
	hosts []string
}

// NewIPHash create new IPHash balancer
func NewIPHash(hosts []string) Balancer {
	return &IPHash{hosts: hosts}
}

// Add new host to the balancer
func (r *IPHash) Add(host string) {
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
func (r *IPHash) Remove(host string) {
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
func (r *IPHash) Balance(key string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(r.hosts))
	return r.hosts[value], nil
}

// Inc .
func (r *IPHash) Inc(_ string) {}

// Done .
func (r *IPHash) Done(_ string) {}
