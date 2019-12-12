package sbt

import (
	"github.com/labbloom/go-merkletree"
)

// generatePresenceProof returns the proof needed for the given indices, the elements for the multiproof, as well as an error. 
func (t *BloomTree) generatePresenceProof(elemIndices []int) (*merkletree.MultiProof, [][]byte, error) {
	data := make([][]byte, len(elemIndices))
	for i, v := range elemIndices {
		for j, vv := range t.state {
			if v == vv[1] {
				data[i] = stringElement(t.state[j][0], t.state[j][1])
				break
			} else if v > vv[1] && v < vv[1]+vv[0] {
				data[i] = stringElement(t.state[j][0], t.state[j][1])
				break
			}
		}
	}
	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}

// generateAbsenceProof returns the proof of absence for given index. To prove the absence, only one
// index is needed. generateAbsenceProof returns the proof, the elements for the multiproof, as well as an error. 
func (t *BloomTree) generateAbsenceProof(index int) (*merkletree.MultiProof, [][]byte, error) {
	var data [][]byte
	var i int
	if i == 0 || i == len(t.state) {
		data = append(data, stringElement(t.state[i][0], t.state[i][1]))
	} else {
		data = append(data, stringElement(t.state[i-1][0], t.state[i+1][1]))
		data = append(data, stringElement(t.state[i+1][0], t.state[i+1][1]))
	}

	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}

// SbtMultiProof returns the multiproof to verify the presence, or absence of an element in a bloom filter.
func (t *BloomTree) SbtMultiProof(elem []byte) (*merkletree.MultiProof, [][]byte, bool, error) {
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