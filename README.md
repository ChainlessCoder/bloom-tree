[![](https://img.shields.io/badge/made%20by-Bloom%20Lab-blue.svg?style=flat-square)](https://bloomlab.io)
[![Build Status](https://travis-ci.com/labbloom/bloom-tree.svg?token=KzkBQ6duyh2GgqS9Be5J&branch=master)](https://travis-ci.com/labbloom/bloom-tree)
[![codecov](https://codecov.io/gh/labbloom/bloom-tree/branch/master/graph/badge.svg?token=xLnQTvQe2W)](https://codecov.io/gh/labbloom/bloom-tree)

# The Bloom Tree
The Bloom Tree combines the idea of bloom filters with that of merkle trees. 
In the standard bloom filter, we are interested to verify the presence, or absence of element(s) in a set. 
In the case of the  Bloom Tree, we are interested to check and transmit the presence, or absence of an element in a secure and bandwidth efficient way to another party. 
Instead of sending the whole bloom filter to a receiver, we only send a small multiproof.

<img src="https://github.com/labbloom/bloom-tree/blob/master/img/bloom_tree.png" class="center" width="700" height="320">

## Install
```sh
go get github.com/labbloom/bloom-tree
```

## Usage
`bloom-tree` generates Merkle tree from `BloomFilter` interface which implements `Proof`, `BitArray`, `MapElementToBF`, `NumOfHashes`, and `GetElementIndicies` (DBF implements all of those methods). Compact multi proofs can be generated and verified. Bloom filter is chunked and then those chunks becomes leaves of merkle tree, the chunks size must be divisible by 64 and can be set using the method SetChunkSize.

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

	// Generate compact multi proof for element "Foo"
	multiproof, err := bt.GenerateCompactMultiProof([]byte("Foo"))
	if err != nil {
		panic(err)
	}
  // if ProofType is equal to 255 then it is presence proof and absence proof otherwise
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
