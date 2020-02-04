package bloomtree

import (
	"crypto/sha512"
	"encoding/binary"
	"errors"
)

var chunkSize = 64

// Hash returns a 256 bit hash
func hashChild(elem1, elem2 [32]byte) [32]byte {
	var elem []byte
	elem = append(elem, elem1[:]...)
	elem = append(elem, elem2[:]...)
	return sha512.Sum512_256(elem)
}

func hashLeaf(index uint64, elements ...uint64) [sha512.Size256]byte {
	var elem []byte

	a := make([]byte, chunkSize)
	binary.LittleEndian.PutUint64(a, index)

	elem = append(elem, a[:]...)
	for _, e := range elements {
		b := make([]byte, 64)
		binary.LittleEndian.PutUint64(b, e)
		elem = append(elem, b...)
	}

	return sha512.Sum512_256(elem)
}

func SetChunkSize(v int) error {
	if v % 64 != 0 {
		return errors.New("The chunk size must be divisible by 64")
	}
	chunkSize = v
	return nil
}
