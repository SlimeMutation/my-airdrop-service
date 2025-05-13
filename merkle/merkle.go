package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type MerkleTree struct {
	Leaves [][]byte
	Layers [][][]byte
}

func NewMerkleTree(leaves [][]byte) *MerkleTree {
	mt := &MerkleTree{Leaves: leaves}
	mt.buildTree()
	return mt
}

func (mt *MerkleTree) buildTree() {
	mt.Layers = append(mt.Layers, mt.Leaves)
	layer := mt.Leaves
	for len(layer) > 1 {
		var nextLayer [][]byte
		for i := 0; i < len(layer); i += 2 {
			left := layer[i]
			var right []byte
			if i+1 < len(layer) {
				right = layer[i+1]
			} else {
				right = left
			}
			hash := sha256.Sum256(append(left, right...))
			nextLayer = append(nextLayer, hash[:])
		}
		mt.Layers = append(mt.Layers, nextLayer)
		layer = nextLayer
	}
}

func (mt *MerkleTree) Root() []byte {
	if len(mt.Layers) == 0 {
		return nil
	}
	return mt.Layers[len(mt.Layers)-1][0]
}

func (mt *MerkleTree) GenerateProof(leafIndex int) [][]byte {
	var proof [][]byte
	index := leafIndex
	for _, layer := range mt.Layers[:len(mt.Layers)-1] {
		pairIndex := index ^ 1
		if pairIndex < len(layer) {
			proof = append(proof, layer[pairIndex])
		}
		index /= 2
	}
	return proof
}

func VerifyProof(leaf []byte, root []byte, proof [][]byte, index int) bool {
	hash := leaf
	for _, sibling := range proof {
		if index%2 == 0 {
			dataMsg := append(hash, sibling...)
			dataHash := sha256.Sum256(dataMsg)
			hash = dataHash[:]
		} else {
			dataMsg := append(sibling, hash...)
			dataHash := sha256.Sum256(dataMsg)
			hash = dataHash[:]
		}
		index /= 2
	}
	fmt.Println("hash===", hex.EncodeToString(hash))
	fmt.Println("root===", hex.EncodeToString(root))
	return bytes.Equal(hash, root)
}
