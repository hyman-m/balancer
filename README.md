# tinybalancer

[![Go Report Card](https://goreportcard.com/badge/github.com/zehuamama/tinybalancer)&nbsp;](https://goreportcard.com/report/github.com/zehuamama/tinybalancer)![GitHub top language](https://img.shields.io/github/languages/top/zehuamama/tinybalancer)&nbsp;![GitHub](https://img.shields.io/github/license/zehuamama/tinybalancer)&nbsp;[![CodeFactor](https://www.codefactor.io/repository/github/zehuamama/tinybalancer/badge)](https://www.codefactor.io/repository/github/zehuamama/tinybalancer)&nbsp;[![codecov](https://codecov.io/gh/zehuamama/tinybalancer/branch/main/graph/badge.svg)](https://codecov.io/gh/zehuamama/tinybalancer)&nbsp; ![go_version](https://img.shields.io/badge/go%20version-1.17-yellow)

tinybalancer is a reverse proxy load balancer that supports http and https. 

* It currently supports four algorithms, namely `round-robin`, `random`, `the power of 2 random choice` and `consistent hash`.
* tinybalancer will periodically perform `healthcheck` on all proxy sites. When the site is unreachable, it will be automatically removed from the balancer. However, tinybalancer will still perform `healthcheck` on unreachable sites. When the site is reachable, it will automatically add it to the balancer.
## Install
First download the source code of balancer:
```shell
> git clone https://github.com/zehuamama/tinybalancer.git
```
compile the source code:
```shell
> cd ./tinybalancer

> go build
```

## Run
Balancer needs to configure the `config.yaml` file, the content is as follows:

```yaml
# Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.


# The load balancing algorithms supported by the balancer are:
# `round-robin` ,`random` ,`p2c`,`consistent-hash`,
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
and now, you can execute `tinybalancer`, the balancer will print the ascii diagram and configuration details:
```shell
> ./tinybalancer

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
## Contributing

If you are intersted in contributing to tinyrpc, please see here: [CONTRIBUTING](https://github.com/zehuamama/tinybalancer/blob/main/CONTRIBUTING.md)

## License

tinybalancer is licensed under the term of the [BSD 2-Clause License](https://github.com/zehuamama/tinybalancer/blob/main/LICENSE)
