// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sync"

	fibHeap "github.com/starwander/GoFibonacciHeap"
)

func init() {
	factories[LeastLoadBalancer] = NewLeastLoad
}

// Tag .
func (h *host) Tag() interface{} { return h.name }

// Key .
func (h *host) Key() float64 { return float64(h.load) }

// LeastLoad will choose a host based on the least load host
type LeastLoad struct {
	sync.RWMutex
	heap *fibHeap.FibHeap
}

// NewLeastLoad create new LeastLoad balancer
func NewLeastLoad(hosts []string) Balancer {
	ll := &LeastLoad{heap: fibHeap.NewFibHeap()}
	for _, h := range hosts {
		ll.Add(h)
	}
	return ll
}

// Add new host to the balancer
func (l *LeastLoad) Add(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok != nil {
		return
	}
	_ = l.heap.InsertValue(&host{hostName, 0})
}

// Remove new host from the balancer
func (l *LeastLoad) Remove(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	_ = l.heap.Delete(hostName)
}

// Balance selects a suitable host according
func (l *LeastLoad) Balance(_ string) (string, error) {
	l.RLock()
	defer l.RUnlock()
	if l.heap.Num() == 0 {
		return "", NoHostError
	}
	return l.heap.MinimumValue().Tag().(string), nil
}

// Inc refers to the number of connections to the server `+1`
func (l *LeastLoad) Inc(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load++
	_ = l.heap.IncreaseKeyValue(h)
}

// Done refers to the number of connections to the server `-1`
func (l *LeastLoad) Done(hostName string) {
	l.Lock()
	defer l.Unlock()
	if ok := l.heap.GetValue(hostName); ok == nil {
		return
	}
	h := l.heap.GetValue(hostName)
	h.(*host).load--
	_ = l.heap.DecreaseKeyValue(h)
}
