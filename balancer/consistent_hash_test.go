// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConsistent_Balance .
func TestConsistent_Balance(t *testing.T) {
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
			NewConsistent([]string{"http://127.0.0.1:1011",
				"http://127.0.0.1:1012", "http://127.0.0.1:1013", "http://127.0.0.1:1014"}),
			"http://127.0.0.1:1011",
			expect{
				"http://127.0.0.1:1012",
				nil,
			},
		},
		{
			"test-2",
			NewConsistent(nil),
			"",
			expect{
				"",
				NoHostError,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			host, err := c.lb.Balance(c.args)
			assert.Equal(t, c.expect.reply, host)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
