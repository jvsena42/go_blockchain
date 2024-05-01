package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jvsena42/go_blockchain/blockchain"
	"github.com/jvsena42/go_blockchain/utils"
	"github.com/jvsena42/go_blockchain/wallet"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainNode struct {
	port uint16
}

func NewBlockchainNode(port uint16) *BlockchainNode {
	return &BlockchainNode{
		port: port,
	}
}

func (bcn *BlockchainNode) Port() uint16 {
	return bcn.port
}

func (bcn *BlockchainNode) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]

	if !ok {
		minerWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minerWallet.BlockchainAddress(), bcn.Port())
		cache["blockchain"] = bc
	}

	return bc
}

func (bcn *BlockchainNode) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		m, _ := bc.MarshalJson()
		io.WriteString(w, string(m[:]))

	default:
		log.Printf("Error: Invalid http request")
	}
}

func (bcn *BlockchainNode) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Add("Content-Type", "application/json")
		bc := bcn.GetBlockchain()
		transactions := bc.TransactionsPool()
		m, _ := json.Marshal(struct {
			Transactions []*blockchain.Transaction `json:"transactions"`
			Length       int                       `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})

		io.WriteString(w, string(m[:]))

	case http.MethodPost:
		decoder := json.NewDecoder((r.Body))
		var t blockchain.TransactionRequest
		err := decoder.Decode(&t)

		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("Error decode")))
			return
		}

		if !t.Valid() {
			log.Println("ERROR: Missing fields!")
			io.WriteString(w, string(utils.JsonStatus("ERROR: Missing fields!")))
			return
		}

		publicKey := utils.StringToPublicKey(*t.SenderPublicKey)
		signature := utils.StringToSignature(*t.Signature)
		bc := bcn.GetBlockchain()

		isCreated := bc.CreateTransaction(*t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, *t.Value, publicKey, signature)

		w.Header().Add("Content-Type", "application/json")
		var responseByte []byte
		if !isCreated {
			w.WriteHeader(http.StatusBadRequest)
			responseByte = utils.JsonStatus("Fail creating transaction")
		} else {
			w.WriteHeader(http.StatusCreated)
			responseByte = utils.JsonStatus("Success!")
		}
		io.WriteString(w, string(responseByte))

	default:
		log.Println("ERROR: Invalid http method")
		w.WriteHeader(http.StatusBadRequest)
	}

}

func (bcn *BlockchainNode) Mine(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcn.GetBlockchain()
		isMined := bc.Mining()

		var messageByte []byte
		if !isMined {
			w.WriteHeader(http.StatusBadRequest)
			messageByte = utils.JsonStatus("Fail mining transacion pool")
		} else {
			messageByte = utils.JsonStatus("Mining success!")
		}
		w.Header().Add("Content-Type", "application/json")
		io.WriteString(w, string(messageByte))

	default:
		log.Println("ERROR: Invalid http method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", bcn.GetChain)
	http.HandleFunc("/transactions", bcn.Transactions)
	http.HandleFunc("/mine", bcn.Mine)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.port)), nil))
}
