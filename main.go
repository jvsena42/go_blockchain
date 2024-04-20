package main

import (
	"fmt"
	"log"

	"github.com/jvsena42/go_blockchain/blockchain"
)

func init() {
	log.SetPrefix("Blockchain node: ")
}

func main() {
	minerAddress := "miner_blockchain_address"

	blockchain := blockchain.NewBlockchain(minerAddress)

	blockchain.AddTransacion("Tony", "Peter", 10089.67897)
	blockchain.AddTransacion("Tony", "Vingadores", 789453123.67897)
	blockchain.Mining()

	blockchain.AddTransacion("Peter", "Pizaria", 0.00000789)
	blockchain.Mining()

	blockchain.AddTransacion("Satoshi Nakamoto", "jvsena42", 89.6697)
	blockchain.Mining()

	blockchain.Print()

	fmt.Println("Balances:")
	fmt.Printf("Miner: %.1f\n", blockchain.CalculateTotalAmount(minerAddress))
	fmt.Printf("Peter: %.1f\n", blockchain.CalculateTotalAmount("Peter"))
	fmt.Printf("jvsena42: %.1f\n", blockchain.CalculateTotalAmount("jvsena42"))
}
