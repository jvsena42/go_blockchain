package main

import (
	"fmt"
	"log"
	"time"
)

type Block struct {
	nonce        int
	previousHash string
	timeStamp    int64
	transactions []string
}

func newBlock(nonce int, previousHash string) *Block {
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

func init() {
	log.SetPrefix("Blockchain node: ")
}

func main() {
	b := newBlock(0, "first hash")
	b.Print()
}
