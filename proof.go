package bloomtree

import (
	"github.com/willf/bitset"
	"math"
	"errors"
	"sort"
	"fmt"
)

type CompactMultiProof struct {
	// Chunks are the leaves of the bloom tree, i.e. the bloom filter values for given parts of the bloom filter.
	Chunks []uint64 
	// Proof are the hashes needed to reconstruct the bloom tree root.
	Proof [][32]byte
}

// newMultiProof generates a Merkle proof
func newCompactMultiProof(chunks []uint64, proof [][32]byte) *CompactMultiProof {
	return &CompactMultiProof{
		Chunks:  chunks,
		Proof: proof,
	}
}

func checkProofType(elemIndices []uint, chunks []uint64) bool {
	if len(elemIndices) == len(chunks) {
		return true
	}
	return false
}

func checkChunkPresence(elemIndices []uint, chunks []uint64) bool {
	for i, v := range elemIndices {
		chunkIndex := uint(math.Floor(float64(v) / float64((chunkSize() ))))
		indexInsideChunk := (v - (chunkIndex * uint(chunkSize()) ))
		chunkBitSet := bitset.From([]uint64{chunks[i]})
		present := chunkBitSet.Test(indexInsideChunk)
		if present != true {
			return false
		}
	}
	return true
}

func computeChunkIndices(elemIndices []uint) []uint64{
	chunkIndices := make([]uint64, len(elemIndices))
	for i, v := range elemIndices {
		index := uint64(math.Floor(float64(v) / float64((chunkSize() ))))
		chunkIndices[i] = uint64(index)
	}
	return chunkIndices
}

func determineOrder2Hash(ind1, indNeighbor int, h1, h2 [32]byte) [32]byte {
	if ind1 > indNeighbor {
		return hashChild(h2,h1)
	}
	return hashChild(h1,h2)
}

func (bt *BloomTree) verifyProof(chunkIndices []uint64, multiproof *CompactMultiProof, root [32]byte) ([][32]byte, error) {
	var (
		pairs []int
		newIndices []uint64
		newBlueNodes [][32]byte
	)
	proof := multiproof.Proof
	blueNodes := make([][32]byte, len(multiproof.Chunks))
	prevIndices := chunkIndices
	indMap := make(map[uint64]int)
	leavesPerLayer := uint64((len(bt.nodes) + 1))
	currentLayer := uint64(0)
	height := int(math.Log2(float64(len(bt.nodes)/2)))
	for i, v := range multiproof.Chunks {
		blueNodes[i] = hashLeaf(v, prevIndices[i])
	}
	for i:=0; i <= height; i ++ {
		if len(newIndices) != 0 {
			for j:=0; j<len(newIndices); j += 2 {
				prevIndices = append(prevIndices, newIndices[j]/2)
			}
			newIndices = nil
		}
		for _, val := range prevIndices {
			neighbor := val^1
			if _, ok := indMap[val+neighbor]; ok {
				if indMap[val+neighbor] != int(val) {
					indMap[val+neighbor] = -1
				}
			} else {
				indMap[val+neighbor] = int(val)
				pairs = append(pairs, int(val+neighbor))
			}
		}
		for k,v := range indMap {
			a,b := order(uint64(v), k - uint64(v))
			newIndices = append(newIndices, a, b)
		}
		sort.Ints(pairs)
		blueNodeNum, proofNum := 0, 0
		for _,v := range pairs {
			value := uint64(v)
			if indMap[value] == -1 {
				newBlueNodes = append(newBlueNodes,hashChild(blueNodes[blueNodeNum], blueNodes[blueNodeNum+1]))
				blueNodeNum += 2
			} else {
				newBlueNodes = append(newBlueNodes,determineOrder2Hash(indMap[value],v-indMap[value], blueNodes[blueNodeNum], proof[proofNum]))
				blueNodeNum ++
				proofNum ++
			}
		}
		blueNodes = newBlueNodes
		newBlueNodes = nil
		blueNodeNum = 0
		indMap = make(map[uint64]int)
		pairs = nil
		leavesPerLayer /= 2
		currentLayer += leavesPerLayer
		prevIndices = nil
	}
	return blueNodes, nil
}



func (bt *BloomTree) VerifyCompactMultiProof(element, seedValue []byte, multiproof *CompactMultiProof, root [32]byte) (bool, error) {
	//var verify bool
	elemIndices := bt.bf.MapElementToBF(element, seedValue)
	chunks := multiproof.Chunks
	chunkIndices := computeChunkIndices(elemIndices)
	if checkProofType(elemIndices, chunks) == true {
		present := checkChunkPresence(elemIndices, chunks)
		if present != true {
			return true, errors.New("The element is not inside the provided chunks for a presence proof")
		}
		verify, err := bt.verifyProof(chunkIndices, multiproof, root)
		if err != nil {
			return true, err
		}
		fmt.Println(verify)
		return true,nil//verify, err
	}
	present := checkChunkPresence(elemIndices, chunks)
	if present == true {
		return false, errors.New("The element cannot be inside the provided chunk for an absence proof")
	}
	verify, err := bt.verifyProof(chunkIndices, multiproof, root)
	if err != nil {
		return false, err
	}
	fmt.Println(verify)
	return false,nil//verify, err
}
