package sbt

import (
	"errors"

	"github.com/labbloom/go-merkletree"
)

// GenerateMultiProof returns the proof needed for given indices, it returns
// the proof, elements in the tree, and error if element is not present in
// the bloom filter, that is if in the corresponding position it has value 0.
func (t *BloomTree) GenerateMultiProof(elemIndices []int) (*merkletree.MultiProof, [][]byte, error) {
	data := make([][]byte, len(elemIndices))
	for i, v := range elemIndices {
		for j, vv := range t.state {
			if v == vv[1] {
				data[i] = stringElement(t.state[j][0], t.state[j][1])
				break
			} else if v < vv[1] {
				if j == 0 {
					return &merkletree.MultiProof{}, nil, errors.New("element is not present")
				}
				if t.state[j-1][0]+t.state[j-1][1] <= v {
					return &merkletree.MultiProof{}, nil, errors.New("element is not present")
				}
				data[i] = stringElement(t.state[j-1][0], t.state[j-1][1])
				break
			}
		}
	}

	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}

// GenerateAbsenceProof returns the proof of absence for given index. To prove the absence only one
// index is needed, it returns the proof, elements in the tree, and error if element maybe is present
// that is if the value of bloom filter in the given index is 1.
func (t *BloomTree) GenerateAbsenceProof(index int) (*merkletree.MultiProof, [][]byte, error) {
	var data [][]byte

	var i int
	for j, vv := range t.state {
		if index == vv[1] {
			return &merkletree.MultiProof{}, nil, errors.New("element can be present")
		} else if index < vv[1] {
			if t.state[j-1][0]+t.state[j-1][1] > index {
				return &merkletree.MultiProof{}, nil, errors.New("element can be present")
			} else {
				j = i - 1
			}
			break
		}
	}

	if i == 0 || i == len(t.state) {
		data = append(data, stringElement(t.state[i][0], t.state[i][1]))
	} else {
		data = append(data, stringElement(t.state[i-1][0], t.state[i+1][1]))
		data = append(data, stringElement(t.state[i+1][0], t.state[i+1][1]))
	}

	proof, err := t.MT.GenerateMultiProof(data)
	return proof, data, err
}
