// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// TestLeastLoad_Balance .
func TestLeastLoad_Balance(t *testing.T) {
	expect, err := Build(LeastLoadBalancer, []string{"192.168.1.1:1015",
		"192.168.1.1:1016", "192.168.1.1:1017", "192.168.1.1:1018"})
	expect.Remove("192.168.1.1:1018")
	assert.Equal(t, err, nil)
	expect.Inc("192.168.1.1:1015")
	expect.Inc("192.168.1.1:1016")
	expect.Inc("192.168.1.1:1016")
	expect.Inc("192.168.1.1:1018")
	expect.Done("192.168.1.1:1018")
	expect.Done("192.168.1.1:1016")
	ll := NewLeastLoad([]string{"192.168.1.1:1016"})
	ll.Remove("192.168.1.1:1018")
	ll.Add("192.168.1.1:1015")
	ll.Add("192.168.1.1:1016")
	ll.Add("192.168.1.1:1017")
	ll.Inc("192.168.1.1:1015")
	ll.Inc("192.168.1.1:1016")
	ll.Inc("192.168.1.1:1016")
	ll.Done("192.168.1.1:1016")
	llHost, _ := ll.Balance("")
	expectHost, _ := expect.Balance("")
	assert.Equal(t, true, reflect.DeepEqual(llHost, expectHost))
}
