package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	TimeStamp    int64
	Nonce        int
	PreviousHash [32]byte
	Transactions []*Transaction
}

func (b *Block) Hash() [32]byte {
	m, _ := b.MarshalJson()
	return sha256.Sum256([]byte(m))
}

func NewBlock(nonce int, previousHash [32]byte, transactions []*Transaction) *Block {
	b := new(Block)
	b.TimeStamp = time.Now().UnixNano()
	b.PreviousHash = previousHash
	b.Nonce = nonce
	b.Transactions = transactions
	return b
}

func (b *Block) Print() {
	fmt.Printf("timestamp:\t%d\n", b.TimeStamp)
	fmt.Printf("nonce:\t\t%d\n", b.Nonce)
	fmt.Printf("previous_hash:\t%x\n", b.PreviousHash)
	for _, t := range b.Transactions {
		t.Print()
	}
}

func (b *Block) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PreviousHash string         `json:"previous_hash"`
		TimeStamp    int64          `json:"time_stamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Nonce:        b.Nonce,
		PreviousHash: fmt.Sprintf("%x", b.PreviousHash),
		TimeStamp:    b.TimeStamp,
		Transactions: b.Transactions,
	})
}
