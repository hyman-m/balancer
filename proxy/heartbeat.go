// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"log"
	"net"
	"time"
)

var (
	HeartBeatTimeout  = 2 * time.Second
	ConnectionTimeout = 1 * time.Second
)

// ReadAlive reads the alive status of the site
func (h *HTTPProxy) ReadAlive(url string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.alive[url]
}

// SetAlive sets the alive status to the site
func (h *HTTPProxy) SetAlive(url string, alive bool) {
	h.Lock()
	defer h.Unlock()
	h.alive[url] = alive
}

// HeartBeat enable a heartbeat goroutine for each agent
func (h *HTTPProxy) HeartBeat() {
	for url := range h.urlMap {
		go h.heartbeat(url)
	}
}

func (h *HTTPProxy) heartbeat(url string) {
	for {
		select {
		case <-time.After(HeartBeatTimeout):
			if !isBackendAlive(url) && h.ReadAlive(url) {
				h.SetAlive(url, false)
				h.lb.Remove(url)
				log.Printf("Site unreachable, remove %s from load balancer.", url)
			} else if isBackendAlive(url) && !h.ReadAlive(url) {
				h.SetAlive(url, true)
				h.lb.Add(url)
				log.Printf("Site reachable, add %s to load balancer.", url)
			}
		}
	}
}

// isBackendAlive Attempt to establish a tcp connection to determine whether the site is alive
func isBackendAlive(url string) bool {
	conn, err := net.DialTimeout("tcp", url, ConnectionTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
