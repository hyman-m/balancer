package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/zehuamama/tinybalancer/balancer"
)

const (
	XRealIP = "X-Real-IP"
	XProxy  = "X-Proxy"
)

var (
	ReverseProxy = "Balancer-Reverse-Proxy"
)

type HTTPProxy struct {
	urlMap map[string]*httputil.ReverseProxy
	lb     balancer.Balancer
}

func NewHTTPProxy(targetHosts []string, algo balancer.Algorithm) (
	*HTTPProxy, error) {
	lb, err := balancer.Build(algo, targetHosts)
	if err != nil {
		return nil, err
	}

	urlMap := make(map[string]*httputil.ReverseProxy)
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

		urlMap[targetHost] = proxy
	}

	return &HTTPProxy{
		urlMap: urlMap,
		lb:     lb,
	}, nil
}

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
