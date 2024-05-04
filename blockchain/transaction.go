package blockchain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	SenderAddress    string
	RecipientAddress string
	Value            float32
}

func (t *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf("sender_blockchain_address:\t%s\n", t.SenderAddress)
	fmt.Printf("recipient_blockchain_address:\t%s\n", t.RecipientAddress)
	fmt.Printf("value:\t\t\t\t%1f\n", t.Value)
}

func (t *Transaction) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress    string  `json:"sender_address"`
		RecipientAddress string  `json:"recipient_address"`
		Value            float32 `json:"value"`
	}{
		SenderAddress:    t.SenderAddress,
		RecipientAddress: t.RecipientAddress,
		Value:            t.Value,
	})
}

func NewTransaction(sender string, recipient string, value float32) *Transaction {

	return &Transaction{
		SenderAddress:    sender,
		RecipientAddress: recipient,
		Value:            value,
	}
}

type TransactionRequest struct {
	SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
	SenderPublicKey            *string  `json:"sender_public_key"`
	Value                      *float32 `json:"value"`
	Signature                  *string  `json:"signature"`
}

func (tr *TransactionRequest) Valid() bool {
	if tr.SenderBlockchainAddress == nil ||
		tr.RecipientBlockchainAddress == nil ||
		tr.SenderPublicKey == nil ||
		tr.Value == nil ||
		tr.Signature == nil {
		return false
	}

	return true
}
