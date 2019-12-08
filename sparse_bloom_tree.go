package sbt

import (
	"fmt"

	"github.com/labbloom/go-merkletree"

	"github.com/willf/bitset"
)

// BloomTree represents sparse bloom tree, field intState
// represents the int state of bloom filter that is the slice of integers
// containing the number of consecutive 1's and the negative value of number
// of consecutive 0's. Field indices the slice of integers containing the
// indices where the bloom filter changes bits. Filed mT represents
// merkle tree build from element combination of two slices of integers.
type BloomTree struct {
	state [][2]int
	mT    *merkletree.MerkleTree
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
					ret = append(ret, [2]int{1, 0})
				} else {
					if b.Test(i - 1) {
						ret[len(ret)-1][0] += 1
					} else {
						ret = append(ret, [2]int{1, int(i)})
					}
				}
			}
		}
	}
	return ret
}

// elementsOfTree prepares elements to build the merkle tree, combination
// of element is done using stringElement function.
func elementsOfTree(state [][2]int) [][]byte {
	length := len(state)

	var elements [][]byte
	for i := 0; i < length; i++ {
		elements = append(elements, stringElement(state[i][0], state[i][1]))

	}

	return elements
}

// merkleTree builds merkle tree with given element
func merkleTree(elements [][]byte) *merkletree.MerkleTree {
	mT, err := merkletree.New(elements)
	if err != nil {
		return &merkletree.MerkleTree{}
	}
	return mT
}

// to build bloom tree is needed only a bloom filter, everything is done
// in this function.
func NewBloomTree(b *bitset.BitSet) *BloomTree {
	state := bit2int(b)
	mT := merkleTree(elementsOfTree(state))

	return &BloomTree{
		mT:    mT,
		state: state,
	}
}

// stringElement returns the combination of two integers
func stringElement(i, j int) []byte {
	return []byte(fmt.Sprintf("%d,%d", i, j))
}
