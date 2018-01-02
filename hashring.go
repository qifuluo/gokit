//Consistent hash
package gokit

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
)

type node struct {
	name     string
	replicas uint32
}

type HashNode struct {
	Hash uint32
	no   *node
}

type HashNodes []*HashNode

func (x HashNodes) Len() int { return len(x) }

func (x HashNodes) Less(i, j int) bool { return x[i].Hash < x[j].Hash }

func (x HashNodes) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type HashRing struct {
	ring HashNodes
}

func NewHashRing() *HashRing {
	return &HashRing{}
}

func (c *HashRing) key(elt string, idx uint32) string {
	return strconv.FormatInt(int64(idx), 10) + elt
}

func (c *HashRing) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}

	return crc32.ChecksumIEEE([]byte(key))
}

func (c *HashRing) Add(elt string, replicas uint32) {
	if 0 == replicas {
		replicas = 1
	}

	no := &node{
		name:     elt,
		replicas: replicas,
	}

	for i := uint32(0); i < replicas; i++ {
		c.ring = append(c.ring, &HashNode{
			Hash: c.hashKey(c.key(elt, i)),
			no:   no,
		})
	}

	sort.Sort(c.ring)
}

// Remove removes an element from the hash.
func (c *HashRing) Remove(elt string) {
	var del uint32
	var replicas uint32
	for i := len(c.ring) - 1; i > 0; i-- {
		if c.ring[i].no.name == elt {
			if 0 == replicas {
				replicas = c.ring[i].no.replicas
			}
			c.ring = append(c.ring[:i], c.ring[i+1:]...)
			del++
			if del == replicas {
				break
			}
		}
	}
}

func (c *HashRing) Get(name string) (string, error) {
	cirLen := c.NodeNum()
	if 0 == cirLen {
		return "", fmt.Errorf("empty circle")
	}
	key := c.hashKey(name)
	return c.ring[c.search(key)].no.name, nil
}

func (c *HashRing) search(key uint32) (i int) {
	f := func(x int) bool {
		return c.ring[x].Hash > key
	}
	i = sort.Search(len(c.ring), f)
	if i >= len(c.ring) {
		i = 0
	}
	return
}

func (c *HashRing) NodeNum() uint32 {
	return uint32(len(c.ring))
}

func (c *HashRing) Print() {
	for i := 0; i < len(c.ring); i++ {
		fmt.Println(fmt.Sprintf("Hash: %v  Name: %v", c.ring[i].Hash, c.ring[i].no.name))
	}
}
