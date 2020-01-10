package bloomtree

import (
	"math"
	"errors"
)

// MerkleProof returns the hashes for a given chunk, up to a specified height.
// index is the index of a bloom filter chunk. intersection is the desired height.
// If intersection is set to 0, MerkleProof returns all the hashes up to the root.
func (bt *BloomTree) MerkleProof(index uint64, intersection int) ([][32]byte, error) {
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