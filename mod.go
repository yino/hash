package hash

import (
	"hash/crc32"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// 一致性 Hash 节点副本数量
const DEFAULT_REPLICAS = 100

// SortKeys 存储一致性 Hash 的值
type SortKeys []uint32

// Len 一致性哈希数量
func (sk SortKeys) Len() int {
	return len(sk)
}

// Less Hash 值的比较
func (sk SortKeys) Less(i, j int) bool {
	return sk[i] < sk[j]
}

// Swap 交换两个 Hash 值
func (sk SortKeys) Swap(i, j int) {
	sk[i], sk[j] = sk[j], sk[i]
}

// Hash 环，存储每一个节点的信息
type HashRing struct {
	Nodes map[uint32]string
	Keys  SortKeys
	sync.RWMutex
}

// New 根据node新建一个Hash环
func (hr *HashRing) New(nodes []string) {
	if nodes == nil {
		return
	}
	hr.Nodes = make(map[uint32]string)
	hr.Keys = SortKeys{}
	for _, node := range nodes {
		// Hash 通过 node 节点名称生成哈希值，该Hash值指向对应 node 节点
		hr.Nodes[hr.Hash(node)] = node
		// 将哈希值保存在key列表中
		hr.Keys = append(hr.Keys, hr.Hash(node))
	}
	// 对 Hash 值进行排序，后面计算的 Hash 值与 Keys 进行比较，取大于等于计算所得的Hash值所对应的node节点
	sort.Sort(hr.Keys)
}

// hashStr 根据Key计算Hash值
func (hr *HashRing) Hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

// GetNode 根据key找出对应的node节点
func (hr *HashRing) GetNode(key string) string {
	hr.RLock()
	defer hr.RUnlock()
	hash := hr.Hash(key)
	i := hr.get_position(hash)
	return hr.Nodes[hr.Keys[i]]
}

// get_position 找出第一个大于等于 hash 的key
func (hr *HashRing) get_position(hash uint32) (i int) {
	i = sort.Search(len(hr.Keys), func(i int) bool {
		return hr.Keys[i] >= hash
	})
	if i >= len(hr.Keys) {
		return 0
	}
	return
}

// NewMod
func NewMod(nodes []string) *HashRing {
	hashR := new(HashRing)
	// 生成 hash 环
	hashR.New(nodes)
	rand.Seed(time.Now().Unix())
	return hashR
}
