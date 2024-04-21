package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"

	"github.com/jvsena42/go_blockchain/utils"
)

type Transaction struct {
	senderPrivateKey *ecdsa.PrivateKey
	senderPublicKey  *ecdsa.PublicKey
	senderAddress    string
	recipientAddress string
	value            float32
}

func NewTransaction(
	privateKey *ecdsa.PrivateKey,
	publicKey *ecdsa.PublicKey,
	sender string,
	recipient string,
	value float32,
) *Transaction {
	return &Transaction{
		senderPrivateKey: privateKey,
		senderPublicKey:  publicKey,
		senderAddress:    sender,
		recipientAddress: recipient,
		value:            value,
	}
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

func (t *Transaction) GenerateSignature() *utils.Signature {
	m, _ := t.MarshalJson()
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, t.senderPrivateKey, h[:])
	return &utils.Signature{R: r, S: s}
}
