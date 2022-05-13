# balancer

[![Go Report Card](https://goreportcard.com/badge/github.com/zehuamama/balancer)&nbsp;](https://goreportcard.com/report/github.com/zehuamama/balancer)![GitHub top language](https://img.shields.io/github/languages/top/zehuamama/balancer)&nbsp;![GitHub](https://img.shields.io/github/license/zehuamama/balancer)&nbsp;[![CodeFactor](https://www.codefactor.io/repository/github/zehuamama/balancer/badge)](https://www.codefactor.io/repository/github/zehuamama/balancer)&nbsp;[![codecov](https://codecov.io/gh/zehuamama/balancer/branch/main/graph/badge.svg)](https://codecov.io/gh/zehuamama/balancer)&nbsp; ![go_version](https://img.shields.io/badge/go%20version-1.17-yellow)

balancer is a reverse proxy load balancer that supports http and https. 

* It currently supports four algorithms, namely `round robin`, `random`, `the power of 2 random choice` , `consistent hash` and `ip hash`.
* `balancer` will perform `health check` on all proxy sites periodically. When the site is unreachable, it will be removed from the balancer automatically . However, `balancer` will still perform `health check` on unreachable sites. When the site is reachable, it will add it to the balancer automatically.
## Install
First download the source code of balancer:
```shell
> git clone https://github.com/zehuamama/balancer.git
```
compile the source code:
```shell
> cd ./balancer

> go build
```

## Run
`Balancer` needs to configure the `config.yaml` file, the content is as follows:

```yaml
# Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


# The load balancing algorithms supported by the balancer are:
# `round-robin` ,`random` ,`p2c`,`consistent-hash`, `ip-hash`
# Among these,`p2c` refers to the power of 2 random choice.

schema: http                      # support http and https
port: 8089                        # port for balancer
ssl_certificate:
ssl_certificate_key:
tcp_health_check: true
health_check_interval: 3          # health check interval (second)
# The maximum number of requests that the balancer can handle at the same time
# 0 refers to no limit to the maximum number of requests
max_allowed: 100
location:                         # route matching for reverse proxy
  - pattern: /
    proxy_pass:                   # URL of the reverse proxy
    - "http://192.168.1.1"
    - "http://192.168.1.2:1015"
    - "https://192.168.1.2"
    - "http://my-server.com"
    balance_mode: round-robin     # load balancing algorithm
```
and now, you can execute `balancer`, the balancer will print the ascii diagram and configuration details:
```shell
> ./balancer

___ _ _  _ _   _ ___  ____ _    ____ _  _ ____ ____ ____ 
 |  | |\ |  \_/  |__] |__| |    |__| |\ | |    |___ |__/ 
 |  | | \|   |   |__] |  | |___ |  | | \| |___ |___ |  \                                        

Schema: http
Port: 8089
Health Check: true
Location:
        Route: /
        Proxy Pass: [http://192.168.1.1 http://192.168.1.2:1015 https://192.168.1.2 http://my-server.com]
        Mode: round-robin

```
## Use API
`balancer` is also a go library that implements load balancing algorithms, it can be used alone as an API, you need to import it into your project first:
```shell
> go get github.com/zehuamama/balancer/balancer
```

Build the load balancer with `balancer.Build`:
```go
hosts := []string{
	"http://192.168.11.101",
	"http://192.168.11.102",
	"http://192.168.11.103",
	"http://192.168.11.104",
}

lb, err := balancer.Build(balancer.P2CBalancer, hosts)
if err != nil {
	return err
}
```
each load balancer implements the `balancer.Balancer` interface:
```go
type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}
```
currently supports the following load balancing algorithms:
```go
const (
	IPHashBalancer         = "ip-hash"
	ConsistentHashBalancer = "consistent-hash"
	P2CBalancer            = "p2c"
	RandomBalancer         = "random"
	R2Balancer             = "round-robin"
)
```
and you can use balancer like this:
```go

clientAddr := "172.160.1.5"  // request IP
	
targetHost, err := lb.Balance(clientAddr) 
if err != nil {
	log.Fatal(err)
}
	
lb.Inc(targetHost)
defer lb.Done(targetHost)

// route to target host
```

## Contributing

If you are intersted in contributing to balancer, please see here: [CONTRIBUTING](https://github.com/zehuamama/balancer/blob/main/CONTRIBUTING.md)

## License

balancer is licensed under the term of the [BSD 2-Clause License](https://github.com/zehuamama/balancer/blob/main/LICENSE)
