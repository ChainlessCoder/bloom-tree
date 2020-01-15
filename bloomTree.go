package bloomtree

import (
	"github.com/willf/bitset"
	"math"
	"errors"
	"sort"
)

// BloomFilter interface. Requires two methods:
// The Bitarray method - returns the bloom filter as a bit array.
// The Proof method - If the element is in the bloom filter, it returns:
// indices, true (where "indices" is an integer array of the indices of the element in the bloom filter).
// If the element is not in the bloom filter, it returns:
// index, false (where "index" is one of the element indices that have a zero value in the bloom filter).
type BloomFilter interface {
	Proof([]byte) ([]uint64, bool)
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
	nodes := make([][32]byte, (leafNum * 2)-1)
	for i,v := range bfAsInt {
		nodes[i] = hashLeaf(v,uint64(i))
	}
	for i:=len(bfAsInt); i < len(nodes) - len(bfAsInt); i ++ {
		nodes[i] = hashPadding(nodes[i], i)
	}
	for i := leafNum-1; i < (leafNum * 2)-1; i++ {
		nodes[i] = hashChild(nodes[i-(leafNum-1)], nodes[i+1-(leafNum-1)])
	} 
	return &BloomTree{
		bf: b,
		nodes: nodes,
	}, nil
}

func order(a,b uint64) (uint64,uint64) {
	if a > b {
		return b,a
	}
	return a,b
	
}


func (bt *BloomTree) generateProof(indices []uint64) ([][32]byte, error) {
	var hashes [][32]byte
	var hashIndices []uint64
	var hashIndicesBucket []int
	var newIndices []uint64
	prevIndices := indices
	indMap := make(map[[2]uint64][2]uint64)
	leavesPerLayer := uint64((len(bt.nodes) + 1))
	currentLayer := uint64(0)
	height := int(math.Log2(float64(len(bt.nodes)/2)))
	for i:=0; i <= height; i ++ {
		if len(newIndices) != 0 {
			for j:=0; j<len(newIndices); j += 2 {
				prevIndices = append(prevIndices, newIndices[j]/2)
			}
			newIndices = nil
		}
		for _, val := range prevIndices {
			neighbor := val^1
			a,b := order(val, neighbor)
			pair := [2]uint64{a,b}
			if _, ok := indMap[pair]; ok {
				indMap[pair] = [2]uint64{1,0}
			} else {
				indMap[pair] = [2]uint64{0, neighbor + currentLayer}
			}
		}
		for k,v := range indMap {
			if v[0] == 0 {
				hashIndicesBucket = append(hashIndicesBucket, int(v[1]))
			}
			newIndices = append(newIndices, k[0], k[1])
		}
		sort.Ints(hashIndicesBucket)
		for _, elem := range hashIndicesBucket {
			hashIndices = append(hashIndices, uint64(elem))
		}
		indMap = make(map[[2]uint64][2]uint64)
		hashIndicesBucket = nil
		leavesPerLayer /= 2
		currentLayer += leavesPerLayer
		prevIndices = nil
	}
	for _, hashInd := range hashIndices {
		hashes = append(hashes, bt.nodes[hashInd])
	}
	return hashes, nil
}

func (bt *BloomTree) getChunksAndIndices(indices []uint64) ([]uint64, []uint64){
	chunks := make([]uint64, len(indices))
	chunkIndices := make([]uint64, len(indices))
	bf := bt.bf.BitArray()
	bfAsInt := bf.Bytes()
	for i, v := range indices {
		index := uint64(math.Floor(float64(v) / float64((chunkSize() * 8))))
		chunks[i] = bfAsInt[index]
		chunkIndices[i] = index
	}
	return chunks, chunkIndices
}


// GenerateCompactMultiProof returns a compact multiproof to verify the presence, or absence of an element in a bloom tree.
func (bt *BloomTree) GenerateCompactMultiProof(elem []byte) (*CompactMultiProof, error) {
	indices, present := bt.bf.Proof(elem)
	chunks, chunkIndices := bt.getChunksAndIndices(indices)
	if present {
		proof, err := bt.generateProof(chunkIndices)
		if err != nil {
			return newCompactMultiProof(nil, nil), err
		}
		return newCompactMultiProof(chunks, proof), nil
	} 
	proof, err := bt.generateProof(indices)
	if err != nil {
		return newCompactMultiProof(nil, nil), err
	}
	return newCompactMultiProof(chunks, proof), nil
}

// Root returns the Bloom Tree root
func (bt *BloomTree) Root() [32]byte {
	return bt.nodes[len(bt.nodes)-1]
}
