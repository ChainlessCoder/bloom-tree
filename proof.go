package bloomtree

import (
	
)

type CompactMultiProof struct {
	// Present is the kind of multiproof. A Presence (true) merkle multiproof, or Absence (false) standard merkle proof 
	Present bool 
	// Chunks are the leaves of the bloom tree, i.e. the bloom filter values for given parts of the bloom filter.
	Chunks []uint64 
	// Proof are the hashes needed to reconstruct the bloom tree root.
	Proof [][32]byte
}

// newMultiProof generates a Merkle proof
func newCompactMultiProof(present bool, chunks []uint64, proof [][32]byte) *CompactMultiProof {
	return &CompactMultiProof{
		Present:  present,
		Chunks:  chunks,
		Proof: proof,
	}
}

/*
func (bt *BloomTree) VerifyCompactMultiProof(multiproof *CompactMultiProof, root [32]byte) (bool, error) {
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
*/

/*
// generatePresenceProof returns the proof needed for the given indices, the elements for the multiproof, as well as an error. 
func (bt *BloomTree) generatePresenceProof(elemIndices []int) (*merkletree.MultiProof, [][]byte, error) {
	chunkIndices := make([]int, len(elemIndices))
	kChunks := make([]uint64, len)(elementIndices)
	hashes := make([][][32]byte, len(elemIndices))
	// Step 1: determine chunks and chunk indices, as well as generate individual proofs
	for i, v := range elemIndices {
		chunkIndices[i] = math.Floor(v / chunkSize())
		kChunks[i] = bt.bf.BitArray().Bytes()[chunkIndices[i]]
		tmpProof, err := t.MerkleProof(chunkIndices[i], 0)
		if err != nil {
			return nil, err
		}
		hashes[i] = tmpProof
	}
	// Step 2: combine the hashes across all proofs and highlight all calculated indices
	



	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}
*/
/*
// generateAbsenceProof returns the proof of absence for given index. To prove the absence, only one
// index is needed. generateAbsenceProof returns the proof, the elements for the multiproof, as well as an error. 
func (t *BloomTree) generateAbsenceProof(index int) (*merkletree.MultiProof, [][]byte, error) {
	var data [][]byte
	if index < t.state[0][1] {
		data = append(data, stringElement(t.state[0][0], t.state[0][1]))
	} else if index > t.state[len(t.state)-1][1] {
		data = append(data, stringElement(t.state[len(t.state)-1][0], t.state[len(t.state)-1][1]))
	} else {
		for i, elm := range t.state {
			if elm[1] > index {
				data = append(data, stringElement(t.state[i-1][0], t.state[i-1][1]))
				data = append(data, stringElement(elm[0], elm[1]))
				break
			}
		}
	}

	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}

// NewBloomTreeProof creates a multi proof for a given element
func NewBloomTreeProof(b BloomFilter, elementIndices []int) []byte {
	indices, present := t.bf.Proof(elem)
	if present {
		proof, data, err := t.generatePresenceProof(indices) 
		if err != nil {
			return &merkletree.MultiProof{}, nil, true, err
		}
		return proof, data, true, err
	} 
	proof, data, err := t.generateAbsenceProof(indices[0])
	if err != nil {
		return &merkletree.MultiProof{}, nil, false, err
	}
	return proof, data, false, err
}


func VerifyProof() {

}

*/