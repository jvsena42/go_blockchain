package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	previousHash [32]byte
	timeStamp    int64
	transactions []*Transaction
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address:\t%s\n", t.senderAddress)
	fmt.Printf("recipient_blockchain_address:\t%s\n", t.recipientAddress)
	fmt.Printf("value:\t\t\t\t%1f\n", t.value)
}

func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.senderAddress,
		RecipientAddress: t.recipientAddress,
		Value:            t.value,
	})
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {

	return &Transaction{
		senderAddress:    sender,
		recipientAddress: recipient,
		value:            value,
	}
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.timeStamp = time.Now().UnixNano()
	b.previousHash = previousHash
	b.nonce = nonce
	b.transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.timeStamp)
	fmt.Printf("nonce:\t\t%d\n", b.nonce)
	fmt.Printf("previous_hash:\t%x\n", b.previousHash)
	for _, t := range b.transactions {
		t.Print()
	}
}

func NewBlockchain() *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	return bc
}

func (b *Block) Hash() [32]byte {
	m, _ := b.MarshalJson()
	return sha256.Sum256([]byte(m))
}

func (b *Block) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash [32]byte       `json:"previous_hash"`
		TimeStamp    int64          `json:"time_stamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.nonce,
		PreviousHash: b.previousHash,
		TimeStamp:    b.timeStamp,
		Transactions: b.transactions,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) AddTransacion(sender string, recipient string, value float32) bool {
	t := NewTransaction(sender, recipient, value)
	bc.transactionPool = append(bc.transactionPool, t)
	return false
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
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
	previousHash := blockchain.LastBlock().previousHash
	blockchain.CreateBlock(10, previousHash)

	blockchain.AddTransacion("Tony", "Peter", 10089.67897)
	blockchain.AddTransacion("Tony", "Vingadores", 789453123.67897)

	previousHash = blockchain.LastBlock().previousHash
	blockchain.CreateBlock(7, previousHash)
	blockchain.AddTransacion("Peter", "Pizaria", 0.00000789)

	previousHash = blockchain.LastBlock().previousHash
	blockchain.CreateBlock(119, previousHash)

	blockchain.AddTransacion("Satoshi Nakamoto", "jvsena42", 89.6697)

	previousHash = blockchain.LastBlock().previousHash
	blockchain.CreateBlock(5871, previousHash)

	blockchain.Print()
}
