// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package balancer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// TestNewIPHash_Add .
func TestNewIPHash_Add(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewIPHash([]string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}),
			"http://127.0.0.1:1012",
			&IPHash{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}},
		},
		{
			"test-2",
			NewIPHash([]string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}),
			"http://127.0.0.1:1013",
			&IPHash{hosts: []string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, true, reflect.DeepEqual(c.expect, c.lb))
		})
	}
}

// TestIPHash_Remove .
func TestIPHash_Remove(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewIPHash([]string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013"}),
			"http://127.0.0.1:1012",
			&IPHash{hosts: []string{"http://127.0.0.1:1011", "http://127.0.0.1:1013"}},
		},
		{
			"test-2",
			NewIPHash([]string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}),
			"http://127.0.0.1:1013",
			&IPHash{hosts: []string{"http://127.0.0.1:1011", "http://127.0.0.1:1012"}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.args)
			assert.Equal(t, true, reflect.DeepEqual(c.expect, c.lb))
		})
	}
}

// TestIPHash_Balance .
func TestIPHash_Balance(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	cases := []struct {
		name   string
		lb     Balancer
		key    string
		expect expect
	}{
		{"test-1",
			NewIPHash([]string{"http://127.0.0.1:1011", "http://127.0.0.1:1012", "http://127.0.0.1:1013"}),
			"192.168.1.1",
			expect{
				"http://127.0.0.1:1011",
				nil,
			},
		},
		{
			"test-2",
			NewIPHash([]string{}),
			"192.168.1.1",
			expect{
				"",
				NoHostError,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			value, err := c.lb.Balance(c.key)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.reply, value))
			assert.Equal(t, true, reflect.DeepEqual(c.expect.err, err))
		})
	}
}
