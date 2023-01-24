// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRandom_Add .
func TestRandom_Add(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012", "http://127.0.0.1:1013"},
				},
				rnd: rnd,
			},
			"http://127.0.0.1:1013",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012", "http://127.0.0.1:1013"},
				},
				rnd: rnd,
			},
		},
		{
			"test-2",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012"},
				},
				rnd: rnd,
			},
			"http://127.0.0.1:1012",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012"},
				},
				rnd: rnd,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRandom_Remove .
func TestRandom_Remove(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012", "http://127.0.0.1:1013"},
				},
				rnd: rnd,
			},
			"http://127.0.0.1:1013",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012"},
				},
				rnd: rnd,
			},
		},
		{
			"test-2",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012"},
				},
				rnd: rnd,
			},
			"http://127.0.0.1:1013",
			&Random{
				BaseBalancer: BaseBalancer{
					hosts: []string{"http://127.0.0.1:1011",
						"http://127.0.0.1:1012"},
				},
				rnd: rnd,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRandom_Balance .
func TestRandom_Balance(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect expect
	}{
		{
			"test-1",
			NewRandom([]string{"http://127.0.0.1:1011"}),
			"",
			expect{
				"http://127.0.0.1:1011",
				nil,
			},
		},
		{
			"test-2",
			NewRandom([]string{}),
			"",
			expect{
				"",
				NoHostError,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply, err := c.lb.Balance(c.args)
			assert.Equal(t, c.expect.reply, reply)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
