package main

import (
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

	address := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	proof, amount, index, err := airdropData.GetProof(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Address: %s\n", address)
	fmt.Printf("index: %d\n", index)
	fmt.Printf("Amount: %d\n", amount)
	fmt.Printf("Proof: %v\n", airdrop.HexProof(proof))

	leafHash := airdrop.MakeLeaf(address, index, amount)

	valid := merkle.VerifyProof(leafHash, root, proof, index)
	if valid {
		fmt.Println("Merkle Proof verified successfully!")
	} else {
		fmt.Println("Merkle Proof verification failed.")
	}
}
