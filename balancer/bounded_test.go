// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// TestBounded_Add .
func TestBounded_Add(t *testing.T) {
	expect, err := Build(BoundedBalancer, []string{"192.168.1.1:1015",
		"192.168.1.1:1016", "192.168.1.1:1017", "192.168.1.1:1018"})
	assert.Equal(t, err, nil)
	bounded := NewBounded(nil)
	bounded.Add("192.168.1.1:1015")
	bounded.Add("192.168.1.1:1016")
	bounded.Add("192.168.1.1:1017")
	bounded.Add("192.168.1.1:1018")
	assert.Equal(t, true, reflect.DeepEqual(expect, bounded))
}

// TestBounded_Remove .
func TestBounded_Remove(t *testing.T) {
	expect, err := Build(BoundedBalancer, []string{"192.168.1.1:1015",
		"192.168.1.1:1016"})
	assert.Equal(t, err, nil)
	bounded := NewBounded([]string{"192.168.1.1:1015",
		"192.168.1.1:1016", "192.168.1.1:1017"})
	bounded.Remove("192.168.1.1:1017")
	assert.Equal(t, true, reflect.DeepEqual(expect, bounded))
}

func TestBounded_Balance(t *testing.T) {
	expect, _ := Build(BoundedBalancer, []string{"192.168.1.1:1015",
		"192.168.1.1:1016", "192.168.1.1:1017", "192.168.1.1:1018"})
	expect.Inc("192.168.1.1:1015")
	expect.Inc("192.168.1.1:1015")
	expect.Inc("NIL")
	expect.Done("192.168.1.1:1015")
	expect.Done("NIL")
	host, _ := expect.Balance("172.166.2.44")
	assert.Equal(t, "192.168.1.1:1017", host)
}
