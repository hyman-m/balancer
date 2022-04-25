package main

import (
	"log"
	"net/http"

	"github.com/zehuamama/tinybalancer/balancer"

	"github.com/zehuamama/tinybalancer/proxy"
)

func main() {
	hosts := []string{"http://127.0.0.1:1012", "http://127.0.0.1:1013"}
	lb, err := balancer.Build("round-robin", hosts)
	if err != nil {
		log.Panic(err)
	}
	p, err := proxy.NewHTTPProxy(hosts, lb)

	if err != nil {
		log.Panic(err)
	}
	http.Handle("/", p)
	http.ListenAndServe(":8089", nil)
}
