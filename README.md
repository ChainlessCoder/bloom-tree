[![](https://img.shields.io/badge/made%20by-Bloom%20Lab-blue.svg?style=flat-square)](https://bloomlab.io)
[![Build Status](https://travis-ci.com/labbloom/bloom-tree.svg?token=KzkBQ6duyh2GgqS9Be5J&branch=master)](https://travis-ci.com/labbloom/bloom-tree)
[![codecov](https://codecov.io/gh/labbloom/bloom-tree/branch/master/graph/badge.svg?token=xLnQTvQe2W)](https://codecov.io/gh/labbloom/bloom-tree)

# The Bloom Tree
The Bloom tree is a probabilistic data structure that combines the idea of Bloom filters with that of Merkle trees. Bloom filters are used to verify the presence, or absence of elements in a set. In the case of the Bloom tree, we are interested to verify and transmit the presence, or absence of an element in a secure and bandwidth efficient way to another party. Instead of sending the whole Bloom filter to a receiver, a compact Merkle multiproof is sent. For more information on the distributed bloom filter, please refer to the original paper.

<img src="https://github.com/labbloom/bloom-tree/blob/master/img/bloom-tree.png" class="center" width="700" height="320">

## Install
```sh
go get github.com/labbloom/bloom-tree
```

## Usage
`bloom-tree` generates a Merkle tree from a `BloomFilter` interface which implements the methods: `Proof`, `BitArray`, `MapElementToBF`, `NumOfHashes`, and `GetElementIndicies` (The [DBF](https://github.com/labbloom/DBF) package implements all of the mentioned methods). To construct a Bloom tree, a given bloom filter gets first split into pre-defined chunks. Those chunks become then leaves of a Merkle tree. The default chunk size is 64 bytes. To change the chunk size, one must use the SetChunkSize method. Chunks must be divisible by 64. 
After construction of the tree, compact Merkle multiproofs can be generated and verified. 


## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/labbloom/DBF"
	bloomtree "github.com/labbloom/bloom-tree"
)

func main() {
	// Data for the bloom filter
	data := [][]byte{
		[]byte("Foo"),
		[]byte("Bar"),
		[]byte("Baz"),
	}
	seed := []byte("secret seed")
	dbf := DBF.NewDbf(200, 0.2, seed)
	for _, d := range data {
		dbf.Add(d)
	}
	// specify the chunk size, the value must be divisible by 64
	bloomtree.SetChunkSize(64)
	// Create the bloom tree
	bt, err := bloomtree.NewBloomTree(dbf)
	if err != nil {
		panic(err)
	}

	// Generate compact multiproof for element "Foo"
	multiproof, err := bt.GenerateCompactMultiProof([]byte("Foo"))
	if err != nil {
		panic(err)
	}
  // if ProofType is equal to 255, it is a presence proof. Any other value means that it is an absence proof.
	if multiproof.ProofType == 255 {
		log.Printf("the proof type for element %s is a presence proof\n", []byte("Foo"))
	} else {
		log.Printf("the proof type for element %s is a absence proof\n", []byte("Foo"))
	}

	verified, err := bloomtree.VerifyCompactMultiProof([]byte("Foo"), seed, multiproof, bt.Root(), bt.GetBloomFilter())
	if err != nil {
		panic(err)
	}
	if !verified {
		panic(fmt.Sprintf("failed to verify proof for %s", []byte("Foo")))
	}
}

```
