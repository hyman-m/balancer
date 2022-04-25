package main

import (
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
		return
	}

	if config.Schema != "http" && config.Schema != "https" {
		log.Fatalf("schema \"%s\" not supported", config.Schema)
		return
	}

	if len(config.Location) == 0 {
		log.Fatalf("location is null")
		return
	}

	router := http.NewServeMux()

	// create proxy
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, balancer.Algorithm(l.BalanceMode))
		if err != nil {
			log.Fatalf("create proxy error: %s", err)
			return
		}
		router.Handle(l.Pattern, httpProxy)
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