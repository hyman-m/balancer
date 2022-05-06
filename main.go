// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/zehuamama/tinybalancer/balancer"
	"github.com/zehuamama/tinybalancer/proxy"
)

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		log.Fatalf("read config error: %s", err)
	}

	err = config.Validation()
	if err != nil {
		log.Fatalf("verify config error: %s", err)
	}

	router := mux.NewRouter()
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, balancer.Algorithm(l.BalanceMode))
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
		}
		// start health check
		if config.HealthCheck {
			httpProxy.HealthCheck(config.HealthCheckInterval)
		}
		router.Handle(l.Pattern, httpProxy)
	}
	if config.MaxAllowed > 0 {
		router.Use(maxAllowedMiddleware(config.MaxAllowed))
	}
	svr := http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	// print config detail
	config.Print()

	// listen and serve
	if config.Schema == "http" {
		err := svr.ListenAndServe()
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	} else if config.Schema == "https" {
		err := svr.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			log.Fatalf("listen and serve error: %s", err)
		}
	}
}

func maxAllowedMiddleware(n uint) mux.MiddlewareFunc {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acquire()
			defer release()
			next.ServeHTTP(w, r)
		})
	}
}
