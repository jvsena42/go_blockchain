package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	MINING_DIFICULTY = 3
	MINING_SENDER    = "BLOCKCHAIN REWARD SYSTEM"
	MINING_REWARD    = 1.0
)

type Block struct {
	timeStamp    int64
	nonce        int
	previousHash [32]byte
	transactions []*Transaction
}

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
}

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockChainAddress string
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientAddress {
				totalAmount += value
			}

			if blockchainAddress == t.senderAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
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

func NewBlockchain(blockChainAddress string) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	bc.blockChainAddress = blockChainAddress
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

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, len(bc.transactionPool))

	for _, t := range bc.transactionPool {
		transactions = append(transactions, NewTransaction(t.senderAddress, t.recipientAddress, t.value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, dificulty int) bool {
	zeros := strings.Repeat("0", dificulty)
	guessBlock := Block{timeStamp: 0, nonce: nonce, previousHash: previousHash, transactions: transactions}
	guessHashString := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashString[:dificulty] == zeros
}

func (bc *Blockchain) ProofOfWOrk() int { //TODO IMPLEMENT GOROUTINES
	transaction := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DIFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) Mining() bool {
	bc.AddTransacion(MINING_SENDER, bc.blockChainAddress, MINING_REWARD)
	previousHash := bc.LastBlock().previousHash
	nonce := bc.ProofOfWOrk()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
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
	minerAddress := "miner_blockchain_address"

	blockchain := NewBlockchain(minerAddress)

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
