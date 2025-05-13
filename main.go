package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/the-web3-contracts/airdrop-service/airdrop"
	"github.com/the-web3-contracts/airdrop-service/merkle"
)

func main() {
	airdropData, err := airdrop.LoadAirdropData("./airdrop.json")
	if err != nil {
		log.Fatal(err)
	}

	root := airdropData.MerkleTree.Root()
	fmt.Printf("Merkle Root: %s\n", hex.EncodeToString(root))

	address := "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	proof, amount, index, err := airdropData.GetProof(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Address: %s\n", address)
	fmt.Printf("Amount: %d\n", amount)
	fmt.Printf("Proof: %v\n", airdrop.HexProof(proof))

	leafData := fmt.Sprintf("%s%d%d", address, index, amount)
	leafHash := sha256.Sum256([]byte(leafData))

	valid := merkle.VerifyProof(leafHash[:], root, proof, index)
	if valid {
		fmt.Println("Merkle Proof verified successfully!")
	} else {
		fmt.Println("Merkle Proof verification failed.")
	}
}
