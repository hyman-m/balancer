# tinybalancer

[![Go Report Card](https://goreportcard.com/badge/github.com/zehuamama/tinybalancer)&nbsp;](https://goreportcard.com/report/github.com/zehuamama/tinybalancer)![GitHub top language](https://img.shields.io/github/languages/top/zehuamama/tinybalancer)&nbsp;![GitHub stars](https://img.shields.io/github/stars/zehuamama/tinybalancer)&nbsp;[![GitHub forks](https://img.shields.io/github/forks/zehuamama/tinybalancer)](https://github.com/zehuamama/tinybalancer/network)&nbsp;![GitHub](https://img.shields.io/github/license/zehuamama/tinybalancer)&nbsp;[![codecov](https://codecov.io/gh/zehuamama/tinybalancer/branch/main/graph/badge.svg)](https://codecov.io/gh/zehuamama/tinybalancer)&nbsp; ![go_version](https://img.shields.io/badge/go%20version-1.17-yellow)

tinybalancer is a reverse proxy load balancer that supports http and https. 

* It currently supports four algorithms, namely `round-robin`, `random`, `the power of 2 random choice` and `consistent hash`.
* tinybalancer will periodically perform heartbeat detection on all proxy sites. When the site is unreachable, it will be automatically removed from the balancer. However, tinybalancer will still perform heartbeat detection on unreachable sites. When the site is reachable, it will automatically add it to the balancer.
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
location:                         # route matching for reverse proxy
  - pattern: /
    proxy_pass:                   # URL of the reverse proxy
    - "http://127.0.0.1:1012"
    - "http://127.0.0.1:1013"
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
Location:
	Route: /
	ProxyPass: [http://127.0.0.1:1012 http://127.0.0.1:1013]
	Mode: round-robin
```
