package bloomtree

// HashFunc256 is a hashing function
type HashFunc256 func(...[]byte) []byte

// Hash returns a 256 bit hash
func (h *HashFunc256) Hash(elem []byte) []byte {
	return sha256.Sum256(elem)
}

// HashLength returns the length of the used hash
func (h *HashFunc256) HashLength(elem []byte) int {
	return 256
}


// HashFunc512 is a hashing function
type HashFunc512 func(...[]byte) []byte

// Hash returns a 512 bit hash
func (h *HashFunc512) Hash(elem []byte) {
	return sha512.Sum512(elem)
}

// HashLength returns the length of the used hash
func (h *HashFunc512) HashLength(elem []byte) int {
	return 512
}


// HashType defines the interface that must be supplied by hash functions
type HashType interface {
	// Hash calculates the hash of a given input
	Hash(...[]byte) []byte

	// HashLength provides the length of the hash
	HashLength() int
}