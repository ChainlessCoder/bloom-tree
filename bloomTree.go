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
	for i:=0; i < len(bfAsInt); i++ {
		nodes[i] = hashLeaf(bfAsInt[i])
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
		index := uint64(math.Ceil(float64((chunkSize() * 8)) / float64(v)) - 1)
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
			return newCompactMultiProof(present, nil, nil), err
		}
		return newCompactMultiProof(present, chunks, proof), nil
	} 
	proof, err := bt.generateProof(indices)
	if err != nil {
		return newCompactMultiProof(present, nil, nil), err
	}
	return newCompactMultiProof(present, chunks, proof), nil
}

/*

// MerkleProof returns the hashes for a given chunk, up to a specified height.
// index is the index of a bloom filter chunk. intersection is the desired height.
// If intersection is set to 0, MerkleProof returns all the hashes up to the root.
func (bt *BloomTree) generateMerkleProof(index uint64, intersection int) ([][32]byte, error) {
	if bt.leafNum() < index {
		return nil, errors.New("index out of range")
	}
	proofLen := bt.height() - intersection
	hashes := make([][32]byte, proofLen)
	cur := 0
	minI := uint64(math.Pow(2, float64(intersection+1))) - 1
	for i := index + uint64(len(bt.nodes)/2); i > minI; i /= 2 {
		hashes[cur] = bt.nodes[i^1]
		cur++
	}
	return hashes, nil
}

// GenerateMultiProof generates the proof for multiple pieces of data.
func (bt *BloomTree) GenerateCompactMultiProof(elementIndices []uint64) ([][32]byte, error) {//(*MultiProof, error) {
	hashes := make([][][32]byte, len(elementIndices))
	indices := make([]uint64, len(elementIndices))

	// Step 1: generate individual proofs
	for i := range elementIndices {
		tmpProof, err := bt.generateMerkleProof(elementIndices[i], 0)
		if err != nil {
			return nil, err
		}
		hashes[i] = tmpProof.Hashes
		indices[i] = tmpProof.Index
	}

	// Step 2: combine the hashes across all proofs and highlight all calculated indices
	proofHashes := make(map[uint64][32]byte)
	calculatedIndices := make([]bool, len(bt.nodes))
	for i, index := range indices {
		hashNum := 0
		for j := uint64(index + uint64(math.Ceil(float64(len(bt.nodes))/2))); j > 1; j /= 2 {
			proofHashes[j^1] = hashes[i][hashNum]
			calculatedIndices[j] = true
			hashNum++
		}
	}

	// Step 3: remove any hashes that can be calculated
	for _, index := range indices {
		for j := uint64(index + uint64(math.Ceil(float64(len(t.nodes))/2))); j > 1; j /= 2 {
			if calculatedIndices[j^1] {
				delete(proofHashes, j^1)
			}
		}
	}
	// 4 prepare compact multiproof order
	compactMultiProof := make([][32]byte, len(proofHashes))
	keysmap := make(map[uint64]uint64)
	keys := make([]int, len(proofHashes))
	height := int(math.Log2(float64(len(bt.nodes)/2)))
	nums := make([]uint64, height)
	n,s := uint64(len(bt.nodes)/2), uint64(0)
	for i := 1; i < height; i ++ {
		nums[height-i-1] = n + s
		s += n
		n /= 2
	}
	nn := 0
	for k := range proofHashes {
		new := math.Exp2(math.Floor(math.Log2(float64(k)))) 
		newind := nums[int((math.Log2(new)))-1] + uint64(math.Abs(float64(new) - float64(k)))
		keys[nn] = int(newind)
		keysmap[newind] = k
		nn ++
	}
	sort.Ints(keys)
	for i, v := range keys {
		compactMultiProof[i] = proofHashes[keysmap[uint64(v)]]
	}	
	return compactMultiProof, nil
}

*/


// Root returns the Bloom Tree root
func (bt *BloomTree) Root() [32]byte {
	return bt.nodes[len(bt.nodes)-1]
}

func (bt *BloomTree) Size() int {
	return len(bt.nodes)
}
