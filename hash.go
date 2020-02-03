package bloomtree

import (
	"crypto/sha512"
	"encoding/binary"
)

// Hash returns a 256 bit hash
func hashChild(elem1, elem2 [32]byte) [32]byte {
	var elem []byte
	elem = append(elem, elem1[:]...)
	elem = append(elem, elem2[:]...)
	return sha512.Sum512_256(elem)
}

func hashLeaf(element, index uint64) [sha512.Size256]byte {
	var elem []byte
	a := make([]byte, chunkSize())
	binary.LittleEndian.PutUint64(a, index)
	b := make([]byte, chunkSize())
	binary.LittleEndian.PutUint64(b, element)
	elem = append(elem, a[:]...)
	elem = append(elem, b[:]...)
	return sha512.Sum512_256(elem)
}

func chunkSize() int {
	return 64
}
