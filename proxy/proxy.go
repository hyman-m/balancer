// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"github.com/zehuamama/tinybalancer/balancer"
	"net"
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
	urlMap map[string]*httputil.ReverseProxy
	lb     balancer.Balancer

	sync.RWMutex // protect alive
	alive        map[string]bool
}

// NewHTTPProxy create  new reverse proxy with url and balancer algorithm
func NewHTTPProxy(targetHosts []string, algo balancer.Algorithm) (
	*HTTPProxy, error) {
	urls := make([]string, 0)

	urlMap := make(map[string]*httputil.ReverseProxy)
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
			req.Header.Set(XRealIP, getIP(req.RemoteAddr))
		}
		alive[url.Host] = true // initial mark alive
		urlMap[url.Host] = proxy
		urls = append(urls, url.Host)
	}

	lb, err := balancer.Build(algo, urls)
	if err != nil {
		return nil, err
	}

	return &HTTPProxy{
		urlMap: urlMap,
		lb:     lb,
		alive:  alive,
	}, nil
}

// ServeHTTP implements a proxy to the http server
func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, err := h.lb.Balance(getIP(r.RemoteAddr))
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		errMsg := fmt.Sprintf("balance error: %s", err.Error())
		w.Write([]byte(errMsg))
		return
	}

	h.lb.Inc(host)
	defer h.lb.Done(host)
	h.urlMap[host].ServeHTTP(w, r)
}

func getIP(remoteAddr string) string {
	remoteHost, _, _ := net.SplitHostPort(remoteAddr)
	return remoteHost
}
