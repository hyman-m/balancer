// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

// TestGetIP .
func TestGetIP(t *testing.T) {
	cases := []struct {
		name   string
		args   string
		expect string
	}{
		{"test-1", "192.168.1.1:80", "192.168.1.1"},
		{"test-2", "171.142.1.0:801", "171.142.1.0"},
		{"test-3", "127.0.0.1:80", "127.0.0.1"},
		{"test-4", "localhost:80", "localhost"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expect, GetIP(c.args))
		})
	}
}

// TestGetHost .
func TestGetHost(t *testing.T) {
	cases := []struct {
		name   string
		args   string
		expect string
	}{
		{"test-1", "http://192.168.1.1:80", "192.168.1.1:80"},
		{"test-2", "https://192.168.1.1:80", "192.168.1.1:80"},
		{"test-3", "https://192.168.1.1", "192.168.1.1:443"},
		{"test-4", "https://test.cn", "test.cn:443"},
		{"test-5", "http://test.cn", "test.cn:80"},
		{"test-6", "https://test.cn:80", "test.cn:80"},
		{"test-7", "wx://test.cn:80", "test.cn:80"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			u, err := url.Parse(c.args)
			assert.Equal(t, nil, err)
			assert.Equal(t, c.expect, GetHost(u))
		})
	}
}

// TestIsBackendAlive .
func TestIsBackendAlive(t *testing.T) {
	cases := []struct {
		name   string
		args   string
		expect bool
	}{
		{"test-1", "www.qq.com:443", true},
		{"test-2", "www.baidu.com:80", true},
		{"test-3", "www.baidu.com:443", true},
		{"test-4", "wt#@!$FS:443", false},
		{"test-4", "wt#@!$FS", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expect, IsBackendAlive(c.args))
		})
	}
}
