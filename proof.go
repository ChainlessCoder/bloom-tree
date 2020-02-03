package bloomtree

import (
	"errors"
	"math"
	"sort"

	"github.com/willf/bitset"
)

type CompactMultiProof struct {
	// Chunks are the leaves of the bloom tree, i.e. the bloom filter values for given parts of the bloom filter.
	Chunks []uint64
	// Proof are the hashes needed to reconstruct the bloom tree root.
	Proof [][32]byte
	// ProofType is 255 if the element is present in the bloom filter. it returns the index of the index if the element is not present in the bloom filter.
	ProofType uint8
}

// newMultiProof generates a Merkle proof
func newCompactMultiProof(chunks []uint64, proof [][32]byte, proofType uint8) *CompactMultiProof {
	return &CompactMultiProof{
		Chunks:    chunks,
		Proof:     proof,
		ProofType: proofType,
	}
}

func CheckProofType(proofType uint8) bool {
	if proofType == maxK {
		return true
	}
	return false
}

func checkChunkPresence(elemIndices []uint, chunks []uint64) bool {
	for i, v := range elemIndices {
		chunkIndex := uint(math.Floor(float64(v) / float64(chunkSize())))
		indexInsideChunk := v - (chunkIndex * uint(chunkSize()))
		chunkBitSet := bitset.From([]uint64{chunks[i]})
		present := chunkBitSet.Test(indexInsideChunk)
		if present != true {
			return false
		}
	}
	return true
}

func computeChunkIndices(elemIndices []uint) []uint64 {
	chunkIndices := make([]uint64, len(elemIndices))
	for i, v := range elemIndices {
		index := uint64(math.Floor(float64(v) / float64(chunkSize())))
		chunkIndices[i] = index
	}
	return chunkIndices
}

func determineOrder2Hash(ind1, indNeighbor int, h1, h2 [32]byte) [32]byte {
	if ind1 > indNeighbor {
		return hashChild(h2, h1)
	}
	return hashChild(h1, h2)
}

func verifyProof(chunkIndices []uint64, multiproof *CompactMultiProof, root [32]byte, treeLength int) (bool, error) {
	var (
		pairs        []int
		newIndices   []uint64
		newBlueNodes [][32]byte
	)

	proof := multiproof.Proof
	blueNodes := make([][32]byte, len(multiproof.Chunks))
	prevIndices := chunkIndices
	indMap := make(map[uint64]int)
	leavesPerLayer := uint64(treeLength + 1)
	currentLayer := uint64(0)
	height := int(math.Log2(float64(treeLength / 2)))
	for i, v := range multiproof.Chunks {
		blueNodes[i] = hashLeaf(v, prevIndices[i])
	}
	// remove duplicates of blue nodes
	var uniqueBlueNodes [][32]byte
	uniqueBlueNodes = append(uniqueBlueNodes, blueNodes[0])
	for i := 1; i < len(blueNodes); i++ {
		if blueNodes[i] != blueNodes[i-1] {
			uniqueBlueNodes = append(uniqueBlueNodes, blueNodes[i])
		}
	}
	blueNodes = uniqueBlueNodes

	// remove duplicates of proof
	var uniqueProof [][32]byte
	uniqueProof = append(uniqueProof, proof[0])
	for i := 1; i < len(proof); i++ {
		if proof[i] != proof[i-1] {
			uniqueProof = append(uniqueProof, proof[i])
		}
	}
	proof = uniqueProof
	proofNum := 0
	for i := 0; i <= height; i++ {
		if len(newIndices) != 0 {
			for j := 0; j < len(newIndices); j += 2 {
				prevIndices = append(prevIndices, newIndices[j]/2)
			}
			newIndices = nil
		}
		for _, val := range prevIndices {
			neighbor := val ^ 1
			if _, ok := indMap[val+neighbor]; ok {
				if indMap[val+neighbor] != int(val) {
					indMap[val+neighbor] = -1
				}
			} else {
				indMap[val+neighbor] = int(val)
				pairs = append(pairs, int(val+neighbor))
			}
		}
		for k, v := range indMap {
			if v == -1 {
				a, b := order((k-1)/2, (k+1)/2)
				newIndices = append(newIndices, a, b)
			} else {
				a, b := order(uint64(v), k-uint64(v))
				newIndices = append(newIndices, a, b)
			}
		}
		sort.Ints(pairs)
		blueNodeNum := 0
		for _, v := range pairs {
			value := uint64(v)
			if indMap[value] == -1 {
				newBlueNodes = append(newBlueNodes, hashChild(blueNodes[blueNodeNum], blueNodes[blueNodeNum+1]))
				blueNodeNum += 2
			} else {
				newBlueNodes = append(newBlueNodes, determineOrder2Hash(indMap[value], v-indMap[value], blueNodes[blueNodeNum], proof[proofNum]))
				blueNodeNum++
				proofNum++
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
	if blueNodes[0] == root {
		return true, nil
	}
	return false, nil
}

// VerifyCompactMultiProof return whether the multi proof provided is true or false.
// The proof type can be absence or presence
func VerifyCompactMultiProof(element, seedValue []byte, multiproof *CompactMultiProof, root [32]byte, bf BloomFilter) (bool, error) {
	// find length of the tree
	dbfBytes := len(bf.BitArray().Bytes())
	if dbfBytes == 0 {
		return false, errors.New("there was no bloom filter provided")
	}
	treeLeafs := int(math.Exp2(math.Ceil(math.Log2(float64(dbfBytes)))))
	treeLength := (treeLeafs * 2) - 1
	elemIndices := bf.MapElementToBF(element, seedValue)
	elemIndicesCopy := elemIndices
	chunks := multiproof.Chunks
	if CheckProofType(multiproof.ProofType) {
		sort.Slice(elemIndices, func(i, j int) bool { return elemIndices[i] < elemIndices[j] })
		chunkIndices := computeChunkIndices(elemIndices)
		present := checkChunkPresence(elemIndices, chunks)
		if present != true {
			return false, errors.New("the element is not inside the provided chunks for a presence proof")
		}
		verify, err := verifyProof(chunkIndices, multiproof, root, treeLength)
		if err != nil {
			return false, err
		}
		return verify, nil //verify, err
	}
	index := []uint{elemIndicesCopy[int(multiproof.ProofType)]}
	chunkIndices := computeChunkIndices(index)
	present := checkChunkPresence(index, chunks)
	if present == true {
		return false, errors.New("the element cannot be inside the provided chunk for an absence proof")
	}
	verify, err := verifyProof(chunkIndices, multiproof, root, treeLength)
	if err != nil {
		return false, err
	}
	return verify, nil //verify, err
}
