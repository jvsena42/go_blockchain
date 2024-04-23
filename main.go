package main

import (
	"fmt"
	"log"

	"github.com/jvsena42/go_blockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain node: ")
}

func main() {
	w := wallet.NewWallet()
	fmt.Println("Private Key: ", w.PrivateKey())
	fmt.Println("Public Key: ", w.PublicKey())
	fmt.Println("BlockchainAddress: ", w.BlockChainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockChainAddress(), "jvsena42", 15684.9)
	fmt.Printf("Signature: %s\n", t.GenerateSignature())
}
