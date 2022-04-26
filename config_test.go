// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReadConfig .
func TestReadConfig(t *testing.T) {
	config, err := ReadConfig("config.yaml")
	expert := &Config{
		SSLCertificateKey: "",
		SSLCertificate:    "",
		Port:              8089,
		Schema:            "http",
		Location: []*Location{
			{
				Pattern: "/",
				ProxyPass: []string{
					"http://127.0.0.1:1012",
					"http://127.0.0.1:1013",
				},
				BalanceMode: "round-robin",
			},
		},
	}
	assert.Equal(t, err, nil)
	assert.Equal(t, config, expert)
}

// TestConfig_Validation .
func TestConfig_Validation(t *testing.T) {
	config, err := ReadConfig("config.yaml")
	assert.Equal(t, err, nil)
	err = config.Validation()
	assert.Equal(t, err, nil)
	config.Print()
}
