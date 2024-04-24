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
	walletPeter := wallet.NewWallet()
	walletJvsena := wallet.NewWallet()

	t := wallet.NewTransaction(walletPeter.PrivateKey(), walletPeter.PublicKey(), walletPeter.BlockChainAddress(), walletJvsena.BlockChainAddress(), 15684.9)
	fmt.Printf("Signature: %s\n", t.GenerateSignature())

	blockchain := blockchain.NewBlockchain(walletMiner.BlockChainAddress())
	isAdded := blockchain.AddTransacion(walletPeter.BlockChainAddress(), walletJvsena.BlockChainAddress(), 15684.9, walletPeter.PublicKey(), t.GenerateSignature())

	fmt.Println("Transaction added to transaction pool? ", isAdded)

	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("Miner has %.1f\n", blockchain.CalculateTotalAmount(walletMiner.BlockChainAddress()))
	fmt.Printf("Jv has %.1f\n", blockchain.CalculateTotalAmount(walletJvsena.BlockChainAddress()))
	fmt.Printf("Peter has %.1f\n", blockchain.CalculateTotalAmount(walletPeter.BlockChainAddress()))
}
