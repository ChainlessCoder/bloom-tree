package sbt

import (
	"fmt"

	"github.com/labbloom/go-merkletree"
	"github.com/willf/bitset"
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

// BloomTree represents the sparse bloom tree (SBT) struct.
// The field "bf" represents the bloom filter used for the SBT.
// The field "state" represents a 2D array, containing the integer state of
// the bloom filter, as well as the corresponding indices for the integer state
// (As shown in the SBT paper).
// The "MT" field represents the merkle tree build on top of the integer state of the bloom filter.
type BloomTree struct {
	state [][2]int
	bf    BloomFilter
	MT    *merkletree.MerkleTree
}

// accept as input bloom filter and returns the slices of integers
// in bloom tree struct
func bit2int(b *bitset.BitSet) [][2]int {
	var ret [][2]int
	last := b.Test(0)
	if last {
		ret = append(ret, [2]int{1, 0})
	}
	length := b.Len()
	if b.Len() > 1 {
		for i := uint(1); i < length; i++ {
			if b.Test(i) {
				if len(ret) == 0 {
					ret = append(ret, [2]int{1, int(i)})
				} else {
					if b.Test(i - 1) {
						ret[len(ret)-1][0] ++
					} else {
						ret = append(ret, [2]int{1, int(i)})
					}
				}
			}
		}
	}
	return ret
}

// elementsOfTree prepares the elements to build the merkle tree, combination
// of element is done using stringElement function.
func elementsOfTree(state [][2]int) [][]byte {
	length := len(state)

	var elements [][]byte
	for i := 0; i < length; i++ {
		elements = append(elements, stringElement(state[i][0], state[i][1]))

	}

	return elements
}

// merkleTree builds merkle tree for a given bloom filter integer array
func merkleTree(elements [][]byte) *merkletree.MerkleTree {
	mT, err := merkletree.New(elements)
	if err != nil {
		return &merkletree.MerkleTree{}
	}
	return mT
}

// NewBloomTree creates a sparse bloom tree.
func NewBloomTree(b BloomFilter) *BloomTree {
	state := bit2int(b.BitArray())
	mT := merkleTree(elementsOfTree(state))

	return &BloomTree{
		MT:    mT,
		state: state,
		bf:    b,
	}
}

// stringElement returns the combination of two integers
func stringElement(i, j int) []byte {
	return []byte(fmt.Sprintf("%d,%d", i, j))
}
