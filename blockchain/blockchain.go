package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jvsena42/go_blockchain/utils"
)

const (
	MINING_DIFICULTY = 3
	MINING_SENDER    = "BLOCKCHAIN REWARD SYSTEM"
	MINING_REWARD    = 1.0
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockChainAddress string
	port              uint16
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

func NewBlockchain(blockChainAddress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.CreateBlock(0, b.Hash())
	bc.blockChainAddress = blockChainAddress
	bc.port = port
	return bc
}

func (bc *Blockchain) TransactionsPool() []*Transaction {
	return bc.transactionPool
}

func (bc *Blockchain) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.chain,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)
	//TODO SYNC NODES
	return isTransacted
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {

		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("Error: not enough balance in wallet")
		// 	return false
		// }

		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Could not verify transaction")
	}
	return false
}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := t.MarshalJson()
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
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
	bc.AddTransaction(MINING_SENDER, bc.blockChainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWOrk()
	previousHash := bc.LastBlock().Hash()
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
