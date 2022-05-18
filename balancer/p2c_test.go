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

// TestP2C_Add .
func TestP2C_Add(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	h1 := &host{"http://127.0.0.1:1011", 0}
	h2 := &host{"http://127.0.0.1:1012", 0}
	h3 := &host{"http://127.0.0.1:1013", 0}
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
			"http://127.0.0.1:1012",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
		},
		{
			"test-2",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
			"http://127.0.0.1:1013",
			&P2C{hosts: []*host{h1, h2, h3}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2, "http://127.0.0.1:1013": h3}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestP2C_Remove .
func TestP2C_Remove(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	h1 := &host{"http://127.0.0.1:1011", 0}
	h2 := &host{"http://127.0.0.1:1012", 0}
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
			"http://127.0.0.1:1012",
			&P2C{hosts: []*host{h1}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1}},
		},
		{
			"test-2",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
			"http://127.0.0.1:1013",
			&P2C{hosts: []*host{h1, h2}, rnd: rnd,
				loadMap: map[string]*host{"http://127.0.0.1:1011": h1,
					"http://127.0.0.1:1012": h2}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestP2C_Balance .
func TestP2C_Balance(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	c := struct {
		name   string
		lb     Balancer
		key    string
		expect expect
	}{
		"test-1",
		NewP2C([]string{"http://127.0.0.1:1011",
			"http://127.0.0.1:1012", "http://127.0.0.1:1013", "http://127.0.0.1:1014"}),
		"key1",
		expect{
			"http://127.0.0.1:1011",
			nil,
		},
	}

	t.Run(c.name, func(t *testing.T) {
		c.lb.Inc("http://127.0.0.1:1013")
		c.lb.Inc("http://127.0.0.1:1013")
		c.lb.Inc("http://127.0.0.1:1")
		c.lb.Done("http://127.0.0.1:1")
		c.lb.Done("http://127.0.0.1:1013")
		host, err := c.lb.Balance(c.key)
		assert.Equal(t, c.expect.reply, host)
		assert.Equal(t, c.expect.err, err)
	})
}
