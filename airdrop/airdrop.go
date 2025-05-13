package airdrop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/the-web3-contracts/airdrop-service/merkle"
)

type AirdropEntry struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type Airdrop struct {
	Entries    []AirdropEntry
	MerkleTree *merkle.MerkleTree
}

func LoadAirdropData(filePath string) (*Airdrop, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var entries []AirdropEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	var leaves [][]byte
	for index, entry := range entries {
		leafData := fmt.Sprintf("%s%d%d", entry.Address, index, entry.Amount)
		leafHash := sha256.Sum256([]byte(leafData))
		leaves = append(leaves, leafHash[:])
	}
	tree := merkle.NewMerkleTree(leaves)
	return &Airdrop{
		Entries:    entries,
		MerkleTree: tree,
	}, nil
}

func (a *Airdrop) GetProof(address string) (proof [][]byte, amount uint64, index int, err error) {
	for i, entry := range a.Entries {
		if entry.Address == address {
			proof = a.MerkleTree.GenerateProof(i)
			return proof, entry.Amount, i, nil
		}
	}
	return nil, 0, 0, fmt.Errorf("address not found")
}

func HexProof(proof [][]byte) []string {
	var hexProof []string
	for _, p := range proof {
		hexProof = append(hexProof, hex.EncodeToString(p))
	}
	return hexProof
}
