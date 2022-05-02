// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"github.com/zehuamama/tinybalancer/balancer"
	"github.com/zehuamama/tinybalancer/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

const (
	XRealIP = "X-Real-IP"
	XProxy  = "X-Proxy"
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

// HTTPProxy refers to a reverse proxy in the balancer
type HTTPProxy struct {
	hostMap map[string]*httputil.ReverseProxy
	lb      balancer.Balancer

	sync.RWMutex // protect alive
	alive        map[string]bool
}

// NewHTTPProxy create  new reverse proxy with url and balancer algorithm
func NewHTTPProxy(targetHosts []string, algo balancer.Algorithm) (
	*HTTPProxy, error) {

	hosts := make([]string, 0)
	hostMap := make(map[string]*httputil.ReverseProxy)
	alive := make(map[string]bool)
	for _, targetHost := range targetHosts {
		url, err := url.Parse(targetHost)
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(url)

		originDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originDirector(req)
			req.Header.Set(XProxy, ReverseProxy)
			if len(req.Header.Get(XRealIP)) == 0 {
				req.Header.Set(XRealIP, util.GetIP(req.RemoteAddr))
			}
		}

		host := util.GetHost(url)
		alive[host] = true // initial mark alive
		hostMap[host] = proxy
		hosts = append(hosts, host)
	}

	lb, err := balancer.Build(algo, hosts)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		hostMap: hostMap,
		lb:      lb,
		alive:   alive,
	}, nil
}

// ServeHTTP implements a proxy to the http server
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("proxy panic :%s", err)
			w.WriteHeader(http.StatusBadGateway)
		}
	}()

	clientIP := util.GetIP(r.RemoteAddr)
	if len(r.Header.Get(XRealIP)) != 0 {
		clientIP = r.Header.Get(XRealIP)
	}

	host, err := h.lb.Balance(clientIP)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		errMsg := fmt.Sprintf("balance error: %s", err.Error())
		w.Write([]byte(errMsg))
		return
	}

	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.hostMap[host].ServeHTTP(w, r)
}
