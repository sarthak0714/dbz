package sharding

import (
	"hash/fnv"
	"sort"
	"strconv"
)

type Ring struct {
	nodes      []string
	vnodeCount int
	hashRing   []uint32
	hashMap    map[uint32]string
}

func NewRing(nodes []string, vnodeCount int) *Ring {
	r := &Ring{
		nodes:      nodes,
		vnodeCount: vnodeCount,
		hashMap:    make(map[uint32]string),
	}
	r.generateHashRing()
	return r
}

func (r *Ring) generateHashRing() {
	for _, node := range r.nodes {
		for i := 0; i < r.vnodeCount; i++ {
			hash := r.hash(node + strconv.Itoa(i))
			r.hashRing = append(r.hashRing, hash)
			r.hashMap[hash] = node
		}
	}
	sort.Slice(r.hashRing, func(i, j int) bool { return r.hashRing[i] < r.hashRing[j] })
}

func (r *Ring) GetNode(key string) string {
	if len(r.hashRing) == 0 {
		return ""
	}

	hash := r.hash(key)
	idx := sort.Search(len(r.hashRing), func(i int) bool { return r.hashRing[i] >= hash })
	if idx == len(r.hashRing) {
		idx = 0
	}
	return r.hashMap[r.hashRing[idx]]
}

func (r *Ring) hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
