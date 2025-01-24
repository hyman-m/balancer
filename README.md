# mini-balancer

[![Go Report Card](https://goreportcard.com/badge/github.com/wanzo-mini/mini-balancer)](https://goreportcard.com/report/github.com/wanzo-mini/mini-balancer)&nbsp;![GitHub top language](https://img.shields.io/github/languages/top/wanzo-mini/mini-balancer)&nbsp;![GitHub](https://img.shields.io/github/license/wanzo-mini/mini-balancer)&nbsp;[![CodeFactor](https://www.codefactor.io/repository/github/wanzo-mini/mini-balancer/badge)](https://www.codefactor.io/repository/github/wanzo-mini/mini-balancer)&nbsp;[![codecov](https://codecov.io/gh/wanzo-mini/mini-balancer/branch/main/graph/badge.svg)](https://codecov.io/gh/wanzo-mini/mini-balancer)&nbsp; ![go_version](https://img.shields.io/badge/go%20version-1.17-yellow)

`mini-balancer` is a layer 7 load mini-balancer that supports http and https, and it is also a go library that implements `load balancing` algorithms.

It currently supports load balancing algorithms: 
* `round-robin`
* `random`
* `power of 2 random choice`
* `consistent hash`
* `consistent hash with bounded`
* `ip-hash`
* `least-load`

## Install
First download the source code of mini-balancer:
```shell
> git clone https://github.com/wanzo-mini/mini-balancer.git
```
compile the source code:
```shell
> cd ./mini-balancer

> go build
```

## Run
`mini-balancer` needs to configure the `config.yaml` file, see [config.yaml](https://github.com/wanzo-mini/mini-balancer/blob/main/config.yaml) :

and now, you can execute `mini-balancer`, the mini-balancer will print the ascii diagram and configuration details:
```shell
> ./mini-balancer

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
`mini-balancer` will perform `health check` on all proxy sites periodically. When the site is unreachable, it will be removed from the mini-balancer automatically . However, `mini-balancer` will still perform `health check` on unreachable sites. When the site is reachable, it will add it to the mini-balancer automatically.

## API Usage
`mini-balancer` is also a go library that implements load balancing algorithms, it can be used alone as an API, you need to import it into your project first:
```shell
> go get github.com/wanzo-mini/mini-balancer/mini-balancer
```

Build the load mini-balancer with `mini-balancer.Build`:
```go
hosts := []string{
	"http://192.168.11.101",
	"http://192.168.11.102",
	"http://192.168.11.103",
	"http://192.168.11.104",
}

lb, err := mini-balancer.Build(mini-balancer.P2Cmini-balancer, hosts)
if err != nil {
	return err
}
```
and you can use mini-balancer like this:
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
each load mini-balancer implements the `mini-balancer.mini-balancer` interface:
```go
type mini-balancer interface {
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
	IPHashmini-balancer         = "ip-hash"
	ConsistentHashmini-balancer = "consistent-hash"
	P2Cmini-balancer            = "p2c"
	Randommini-balancer         = "random"
	R2mini-balancer             = "round-robin"
	LeastLoadmini-balancer      = "least-load"
	Boundedmini-balancer        = "bounded"
)
```


## Contributing

If you are interested in contributing to mini-balancer, please see here: [CONTRIBUTING](https://github.com/wanzo-mini/mini-balancer/blob/main/CONTRIBUTING.md)

## License

mini-balancer is licensed under the term of the [BSD 2-Clause License](https://github.com/wanzo-mini/mini-balancer/blob/main/LICENSE)
