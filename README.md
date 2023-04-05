# balancer

[![Go Report Card](https://goreportcard.com/badge/github.com/zehuamama/balancer)&nbsp;](https://goreportcard.com/report/github.com/zehuamama/balancer)![GitHub top language](https://img.shields.io/github/languages/top/zehuamama/balancer)&nbsp;![GitHub](https://img.shields.io/github/license/zehuamama/balancer)&nbsp;[![CodeFactor](https://www.codefactor.io/repository/github/zehuamama/balancer/badge)](https://www.codefactor.io/repository/github/zehuamama/balancer)&nbsp;[![codecov](https://codecov.io/gh/zehuamama/balancer/branch/main/graph/badge.svg)](https://codecov.io/gh/zehuamama/balancer)&nbsp; ![go_version](https://img.shields.io/badge/go%20version-1.17-yellow)

`balancer` is a layer 7 load balancer that supports http and https, and it is also a go library that implements `load balancing` algorithms.

It currently supports load balancing algorithms: 
* `round-robin`
* `random`
* `power of 2 random choice`
* `consistent hash`
* `consistent hash with bounded`
* `ip-hash`
* `least-load`

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
`Balancer` needs to configure the `config.yaml` file, see [config.yaml](https://github.com/zehuamama/balancer/blob/main/config.yaml) :

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
`balancer` will perform `health check` on all proxy sites periodically. When the site is unreachable, it will be removed from the balancer automatically . However, `balancer` will still perform `health check` on unreachable sites. When the site is reachable, it will add it to the balancer automatically.

## API Usage
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
	LeastLoadBalancer      = "least-load"
	BoundedBalancer        = "bounded"
)
```


## Contributing

If you are interested in contributing to balancer, please see here: [CONTRIBUTING](https://github.com/zehuamama/balancer/blob/main/CONTRIBUTING.md)

## License

balancer is licensed under the term of the [BSD 2-Clause License](https://github.com/zehuamama/balancer/blob/main/LICENSE)
