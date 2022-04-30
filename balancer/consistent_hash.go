// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"sort"
	"sync"
	"time"
)

//节点副本数量
const replicationFactor = 10
type Int32SortSet []uint32

// Search 查找第一个大于等于hostHash的哈希值（查找第一个大于等于hostHash的节点）
func (sortedSet Int32SortSet) Search(hostHash uint32) uint32 {
	//如果比最大的哈希值大，那么返回最小的
	if hostHash > sortedSet[len(sortedSet)-1] || len(sortedSet) == 1{
		return sortedSet[0]
	}
	if hostHash < sortedSet[0]{
		return sortedSet[0]
	}
	index:=sort.Search(len(sortedSet) ,func(i int) bool { return sortedSet[i]>= hostHash })
	return sortedSet[index]
}

//PUSH 插入元素同时维护有序集合
func (sortedSet *Int32SortSet)PUSH(hostHash uint32)  {
	if len(*sortedSet) == 0{
		*sortedSet = append(*sortedSet,hostHash)
		return
	}

	var mid int
	low := 0
	hig := len(*sortedSet)-1
	var suiteIndex int
	//如果比最大的哈希值大，那么返回最小的哈希值
	for low < hig{
		mid = low + (hig-low)/2
		if (*sortedSet)[mid] < hostHash {
			low = mid + 1
		}else if  (*sortedSet)[mid] > hostHash {
			hig = mid - 1
		}else { //在set里有此元素故不append
			return
		}
	}
	suiteIndex = low
	if (*sortedSet)[low] < hostHash {
		suiteIndex = low + 1
	}
	//新建临时切片
	temp :=  append(make(Int32SortSet,0),(*sortedSet)[suiteIndex:]...)
	*sortedSet = append((*sortedSet)[:suiteIndex], hostHash)
	*sortedSet = append(*sortedSet,temp...)
	return
}
//Remove 删除指定元素
func (sortedSet *Int32SortSet)Remove(hostHash uint32){
	index := sort.Search(len(*sortedSet) ,func(i int) bool { return (*sortedSet)[i]==hostHash})
	*sortedSet = append((*sortedSet)[:index],(*sortedSet)[index+1:]...)
}



func init() {
	factories["consistent-hash"] = NewConsistent
}

// Consistent refers to consistent hash
type Consistent struct {
	sync.RWMutex
	rnd   *rand.Rand
	hosts     map[uint32]string
	sortedSet *Int32SortSet
}

// NewConsistent create new Consistent balancer
func NewConsistent(hosts []string) Balancer {
	p:=new(Int32SortSet)
	*p = make([]uint32,0)
	c := &Consistent{sortedSet: p,
		hosts : make(map[uint32]string),
		rnd: rand.New(rand.NewSource(time.Now().UnixNano()))}
	for _, host :=range hosts{
		c.Add(host)
	}
	return c
}

// Add new host to the balancer
func (c *Consistent) Add(host string) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < replicationFactor; i++ {
		var hashHost uint32
		if i == 0{
			hashHost = c.hash(host)
		}else {
			hashHost = c.hash(fmt.Sprintf("%s%d", host, i))
		}
		c.hosts[hashHost] = host
		c.sortedSet.PUSH(hashHost)
	}
}

// Remove new host from the balancer
func (c *Consistent) Remove(host string) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < replicationFactor; i++ {
		hashHost := c.hash(fmt.Sprintf("%s%d", host, i))
		delete(c.hosts, hashHost)
		c.sortedSet.Remove(hashHost)
	}
}

// Balance selects a suitable host according to the key value
func (c *Consistent) Balance(key string) (string, error) {
	if len(c.hosts) == 0 {
		return "", NoHostError
	}
	//对请求key进行hash
	hashHost := c.hash(key)
	a:=c.sortedSet.Search(hashHost)
	return c.hosts[a],nil
}

// Inc .
func (c *Consistent) Inc(_ string) {}

// Done .
func (c *Consistent) Done(_ string) {}

func (c *Consistent) hash(key string) (uint32) {
	return crc32.ChecksumIEEE([]byte(key))
}

