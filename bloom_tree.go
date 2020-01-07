package bloomtree

import (
	"github.com/willf/bitset"
	"math"
	"errors"
)


// BloomFilter interface. Requires two methods:
// The Bitarray method - returns the bloom filter as a bit array.
// The Proof method - If the element is in the bloom filter, it returns:
// indices, true (where "indices" is an integer array of the indices of the element in the bloom filter).
// If the element is not in the bloom filter, it returns:
// index, false (where "index" is one of the element indices that have a zero value in the bloom filter).
type BloomFilter interface {
	Proof([]byte) ([]int, bool)
	BitArray() *bitset.BitSet
}

// BloomTree represents the bloom tree struct.
type BloomTree struct {
	bf BloomFilter
	nodes [][32]byte
}

func treeNodes(l []uint64, leafNum int) [][32]byte {
	nodes := make([][32]byte, (leafNum * 2) - 1)
	for i:=0; i < len(l); i++ {
		nodes[i] = hashLeaf(l[i])
	}
	return nodes
}


// NewBloomTree creates a new bloom tree.
func NewBloomTree(b BloomFilter) (*BloomTree, error) {
	bf := b.BitArray()
	bfAsInt := bf.Bytes() 
	if len(bfAsInt) == 0 {
		return nil, errors.New("tree must have at least 1 element")
	}
	leafNum := int(math.Exp2(math.Ceil(math.Log2(float64(len(bfAsInt))))))
	nodes := treeNodes(bfAsInt, leafNum)
	num, num1 := leafNum, 0
	for j:=2; j <= leafNum; j *= 2 {
		for i:=0; i < leafNum / j; i+=2 {
			nodes[num+i] = hashChild(nodes[num1+i], nodes[num1+i+1])
		}
		num += leafNum / j
		num1 += leafNum / (j/2)
	}
	return &BloomTree{
		bf: b,
		nodes: nodes,
	}, nil
}

// Root returns the Bloom Tree root
func (bt *BloomTree) Root() [32]byte {
	bloomFilter := bt.bf.BitArray()
	bfAsInt := bloomFilter.Bytes() 
	nodesArraySize := (int(math.Exp2(math.Ceil(math.Log2(float64(len(bfAsInt)))))) * 2) - 1
	return bt.nodes[nodesArraySize - 1]
}
/*
func NewMerkleProof(b BloomFilter, elementIndices []int) []byte {
	return []byte
}
*/