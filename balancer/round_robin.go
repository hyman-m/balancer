package balancer

import (
	"errors"
	"sync"
)

type RoundRobin struct {
	i     uint64
	hosts []string
	sync.Mutex
}

func init() {
	factories["round-robin"] = NewRoundRobin
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{i: 0, hosts: hosts}
}

func (r *RoundRobin) Add(host string) {
	r.Lock()
	defer r.Unlock()
	for _, h := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts = append(r.hosts, host)
}

func (r *RoundRobin) Remove(host string) bool {
	r.Lock()
	defer r.Unlock()
	for i, h := range r.hosts {
		if h == host {
			r.hosts = append(r.hosts[:i], r.hosts[i+1:]...)
			return true
		}
	}
	return false
}

func (r *RoundRobin) Balance(_ string) (string, error) {
	r.Lock()
	defer r.Unlock()
	if len(r.hosts) == 0 {
		return "", errors.New("no host")
	}
	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.i++
	return host, nil
}

func (r *RoundRobin) Inc(_ string) {}

func (r *RoundRobin) Done(_ string) {}
