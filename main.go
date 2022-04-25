package main

import (
	"net/http"
	"strconv"

	"github.com/fatih/color"
	"github.com/zehuamama/tinybalancer/balancer"
	"github.com/zehuamama/tinybalancer/proxy"
)

func main() {
	config, err := ReadConfig("config.yaml")
	if err != nil {
		color.Red("read config error: %s", err)
		return
	}

	if config.Schema != "http" && config.Schema != "https" {
		color.Red("schema \"%s\" not supported", config.Schema)
		return
	}

	if len(config.Location) == 0 {
		color.Red("location is null")
		return
	}

	router := http.NewServeMux()

	// create proxy
	for _, l := range config.Location {
		httpProxy, err := proxy.NewHTTPProxy(l.ProxyPass, balancer.Algorithm(l.BalanceMode))
		if err != nil {
			color.Red("create proxy error: %s", err)
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
			color.Red("listen and serve error: %s", err)
		}
	} else if config.Schema == "https" {
		err := svr.ListenAndServeTLS(config.SSLCertificate, config.SSLCertificateKey)
		if err != nil {
			color.Red("listen and serve error: %s", err)
		}
	}
}
