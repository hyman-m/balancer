// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"fmt"
	"net"
	"net/url"
	"time"
)

var ConnectionTimeout = 2 * time.Second

func GetIP(remoteAddr string) string {
	remoteHost, _, _ := net.SplitHostPort(remoteAddr)
	return remoteHost
}

func GetHost(url *url.URL) string {
	if _, _, err := net.SplitHostPort(url.Host); err == nil {
		return url.Host
	}
	if url.Scheme == "http" {
		return fmt.Sprintf("%s:%s", url.Host, "80")
	} else if url.Scheme == "https" {
		return fmt.Sprintf("%s:%s", url.Host, "443")
	}
	return url.Host
}

// IsBackendAlive Attempt to establish a tcp connection to determine whether the site is alive
func IsBackendAlive(host string) bool {
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return false
	}
	resolve := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	conn, err := net.DialTimeout("tcp", resolve, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
