package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/jvsena42/go_blockchain/utils"
)

const (
	MINING_DIFICULTY = 3
	MINING_SENDER    = "BLOCKCHAIN REWARD SYSTEM"
	MINING_REWARD    = 1.0
	MINING_TIMER_SEC = 20

	BLOCKCHAIN_PORT_RANGE_START       = 3333
	BLOCKCHAIN_PORT_RANGE_END         = 3336
	NEIGHBOR_IP_RANGE_START           = 0
	NEIGHBOR_IP_RANGE_END             = 3
	BLOCKCHAIN_NEIGBHOR_SYNC_TIME_SEC = 10
)

type Blockchain struct {
	TransactionPool   []*Transaction
	Chain             []*Block
	BlockChainAddress string
	Port              uint16
	mux               sync.Mutex

	neighbors    []string
	muxNeighbors sync.Mutex
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			value := t.Value
			if blockchainAddress == t.RecipientAddress {
				totalAmount += value
			}

			if blockchainAddress == t.SenderAddress {
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
	bc.BlockChainAddress = blockChainAddress
	bc.Port = port
	return bc
}

func (bc *Blockchain) Run() {
	bc.StartSyncNeighbors()
}

func (bc *Blockchain) SetNeightbors() {
	bc.neighbors = utils.FindNeighbors(
		utils.GetHost(),
		bc.Port,
		NEIGHBOR_IP_RANGE_START,
		NEIGHBOR_IP_RANGE_END,
		BLOCKCHAIN_PORT_RANGE_START,
		BLOCKCHAIN_PORT_RANGE_END,
	)

	if len(bc.neighbors) > 0 {
		log.Println("This node's neighbors are", bc.neighbors)
	} else {
		log.Println("This node could not find neighbors", bc.neighbors)
	}
}

func (bc *Blockchain) SyncNeighbors() {
	bc.muxNeighbors.Lock()
	defer bc.muxNeighbors.Unlock()
	bc.SetNeightbors()
}

func (bc *Blockchain) StartSyncNeighbors() {
	bc.SyncNeighbors()
	_ = time.AfterFunc(time.Second*BLOCKCHAIN_NEIGBHOR_SYNC_TIME_SEC, bc.StartSyncNeighbors)
}

func (bc *Blockchain) TransactionsPool() []*Transaction {
	return bc.TransactionPool
}

func (bc *Blockchain) ClearTransactionPool() {
	bc.TransactionPool = bc.TransactionPool[:0]
}

func (bc *Blockchain) MarshalJson() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{
		Blocks: bc.Chain,
	})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.TransactionPool)
	bc.Chain = append(bc.Chain, b)
	bc.TransactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) CreateTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	isTransacted := bc.AddTransaction(sender, recipient, value, senderPublicKey, s)

	if isTransacted {
		bc.broadcasTransaction(senderPublicKey, s, sender, recipient, value)
	}

	return isTransacted
}

func (bc *Blockchain) broadcasTransaction(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, sender string, recipient string, value float32) {
	for _, neighborIPAddress := range bc.neighbors {
		publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
		signatureStr := s.String()
		bt := &TransactionRequest{
			&sender,
			&recipient,
			&publicKeyStr,
			&value,
			&signatureStr}
		marshalJson, _ := json.Marshal(bt)

		buf := bytes.NewBuffer(marshalJson)
		endpoint := fmt.Sprintf("http://%s/transactions", neighborIPAddress)

		client := &http.Client{}
		request, _ := http.NewRequest("PUT", endpoint, buf)
		response, err := client.Do(request)
		if err != nil {
			log.Printf("%v", response)
		}
	}
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32, senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.TransactionPool = append(bc.TransactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {

		// if bc.CalculateTotalAmount(sender) < value {
		// 	log.Println("Error: not enough balance in wallet")
		// 	return false
		// }

		bc.TransactionPool = append(bc.TransactionPool, t)
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
	transactions := make([]*Transaction, len(bc.TransactionPool))

	for _, t := range bc.TransactionPool {
		transactions = append(transactions, NewTransaction(t.SenderAddress, t.RecipientAddress, t.Value))
	}

	return transactions
}

func (bc *Blockchain) ValidProof(nonce int, previousHash [32]byte, transactions []*Transaction, dificulty int) bool {
	zeros := strings.Repeat("0", dificulty)
	guessBlock := Block{TimeStamp: 0, Nonce: nonce, PreviousHash: previousHash, Transactions: transactions}
	guessHashString := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashString[:dificulty] == zeros
}

func (bc *Blockchain) ProofOfWOrk() int {
	transaction := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	for !bc.ValidProof(nonce, previousHash, transaction, MINING_DIFICULTY) {
		nonce++
	}
	return nonce
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *Blockchain) StartMining() {
	bc.Mining()
	_ = time.AfterFunc(time.Second*MINING_TIMER_SEC, bc.StartMining)
}

func (bc *Blockchain) Mining() bool {
	bc.mux.Lock()
	defer bc.mux.Unlock()

	if len(bc.TransactionPool) == 0 {
		return false
	}

	bc.AddTransaction(MINING_SENDER, bc.BlockChainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWOrk()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Chain {
		fmt.Printf("%s Block %d %s\n", strings.Repeat("=", 10), i, strings.Repeat("=", 10))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("#", 30))
}

type AmountResponse struct {
	Amount float32 `json:"amount"`
}
