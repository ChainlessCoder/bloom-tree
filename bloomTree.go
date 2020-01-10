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


// NewBloomTree creates a new bloom tree.
func NewBloomTree(b BloomFilter) (*BloomTree, error) {
	bf := b.BitArray()
	bfAsInt := bf.Bytes() 
	if len(bfAsInt) == 0 {
		return nil, errors.New("tree must have at least 1 element")
	}
	leafNum := int(math.Exp2(math.Ceil(math.Log2(float64(len(bfAsInt))))))
	nodes := make([][32]byte, leafNum+len(bfAsInt)+(leafNum-len(bfAsInt)))
	for i:=0; i < len(bfAsInt); i++ {
		nodes[i+leafNum] = hashLeaf(bfAsInt[i])
	}
	for i := leafNum - 1; i > 0; i-- {
		nodes[i] = hashChild(nodes[i*2], nodes[i*2+1])
	} 
	return &BloomTree{
		bf: b,
		nodes: nodes,
	}, nil
}

// Root returns the Bloom Tree root
func (bt *BloomTree) Root() [32]byte {
	return bt.nodes[1]
}

func (bt *BloomTree) height() int {
	bf := bt.bf.BitArray()
	return int(math.Log2(math.Exp2(math.Ceil(math.Log2(float64(len(bf.Bytes())))))))
}

func (bt *BloomTree) leafNum() uint64 {
	bf := bt.bf.BitArray()
	return uint64(len(bf.Bytes()))
}