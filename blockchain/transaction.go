package blockchain

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	senderAddress    string
	recipientAddress string
	value            float32
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
