// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"hash/crc32"
)

func init() {
	factories[IPHashBalancer] = NewIPHash
}

// IPHash will choose a host based on the client's IP address
type IPHash struct {
	BaseBalancer
}

// NewIPHash create new IPHash balancer
func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
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
