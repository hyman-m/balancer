// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ConnectionTimeout refers to connection timeout for health check
var ConnectionTimeout = 3 * time.Second

// GetIP get client IP
func GetIP(r *http.Request) string {
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	if len(r.Header.Get(XForwardedFor)) != 0 {
		xff := r.Header.Get(XForwardedFor)
		s := strings.Index(xff, ", ")
		if s == -1 {
			s = len(r.Header.Get(XForwardedFor))
		}
		clientIP = xff[:s]
	} else if len(r.Header.Get(XRealIP)) != 0 {
		clientIP = r.Header.Get(XRealIP)
	}

	return clientIP
}

// GetHost get the hostname, looks like IP:Port
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
	resolveAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
	conn, err := net.DialTimeout("tcp", resolveAddr, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
