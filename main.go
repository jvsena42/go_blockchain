package main

import (
	"fmt"
	"log"

	"github.com/jvsena42/go_blockchain/blockchain"
	"github.com/jvsena42/go_blockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain node: ")
}

func main() {
	walletMiner := wallet.NewWallet()
	walletAlice := wallet.NewWallet()
	walletBob := wallet.NewWallet()

	// wallet transaction request
	t := wallet.NewTransaction(walletAlice.PrivateKey(), walletAlice.PublicKey(), walletAlice.BlockchainAddress(), walletBob.BlockchainAddress(), 23.0)

	// blockchain node transaction request handling
	blockchain := blockchain.NewBlockchain(walletMiner.BlockchainAddress())
	isAdded := blockchain.AddTransaction(walletAlice.BlockchainAddress(), walletBob.BlockchainAddress(), 23.0, walletAlice.PublicKey(), t.GenerateSignature())
	fmt.Println("Transaction added to transaction pool?", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("Miner has %.1f\n", blockchain.CalculateTotalAmount(walletMiner.BlockchainAddress()))
	fmt.Printf("Alice has %.1f\n", blockchain.CalculateTotalAmount(walletAlice.BlockchainAddress()))
	fmt.Printf("Bob has %.1f\n", blockchain.CalculateTotalAmount(walletBob.BlockchainAddress()))
}
