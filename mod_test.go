package hash_test

import (
	"fmt"
	"testing"

	"github.com/yino/hash"
)

// TestHashMod  hash取余
func TestHashMod(t *testing.T) {
	var nodes []string
	// 节点数
	nodeREPLICAS := 100
	// 初始化 node 节点
	for i := 1; i <= nodeREPLICAS; i++ {
		nodes = append(nodes, fmt.Sprintf("test_%d", i))
	}
	hashR := hash.NewMod(nodes)
	fmt.Println(hashR.GetNode("test 1"))
}
