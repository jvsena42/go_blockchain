package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timeStamp    int64
	transactions []string
}

func NewBlock(nonce int, previousHash string) *Block {
	b := new(Block)
	b.timeStamp = time.Now().UnixNano()
	b.previousHash = previousHash
	b.nonce = nonce
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.timeStamp)
	fmt.Printf("nonce:\t\t%d\n", b.nonce)
	fmt.Printf("previous_hash:\t%s\n", b.previousHash)
	fmt.Printf("transactions:\t%s\n", b.transactions)
}

type Blockchain struct {
	transactionPool []string
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	bc := new(Blockchain)
	bc.CreateBlock(0, "hash #0 genesis block")
	return bc
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash string) *Block {
	b := NewBlock(nonce, previousHash)
	bc.chain = append(bc.chain, b)
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 10), i, strings.Repeat("=", 10))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("#", 30))
}

func init() {
	log.SetPrefix("Blockchain node: ")
}

func main() {
	blockchain := NewBlockchain()

	blockchain.CreateBlock(89, "hash #1")

	blockchain.CreateBlock(5, "hash #2")

	blockchain.CreateBlock(32, "hash #3")

	blockchain.CreateBlock(44, "hash #4")
	blockchain.Print()
}
