package airdrop

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
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

func MakeLeaf(address string, index int, amount uint64) []byte {
	addrBytes, _ := hex.DecodeString(address[2:]) // 去掉0x
	idxBytes := make([]byte, 32)
	amtBytes := make([]byte, 32)
	big.NewInt(int64(index)).FillBytes(idxBytes)
	big.NewInt(int64(amount)).FillBytes(amtBytes)
	packed := append(addrBytes, append(idxBytes, amtBytes...)...)
	hash := sha256.Sum256(packed)
	return hash[:]
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
		leafHash := MakeLeaf(entry.Address, index, entry.Amount)
		leaves = append(leaves, leafHash)
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
