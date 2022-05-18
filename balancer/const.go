// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

const (
	IPHashBalancer         = "ip-hash"
	ConsistentHashBalancer = "consistent-hash"
	P2CBalancer            = "p2c"
	RandomBalancer         = "random"
	R2Balancer             = "round-robin"
	LeastLoadBalancer      = "least-load"
	BoundedBalancer        = "bounded"
)
