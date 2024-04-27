package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jvsena42/go_blockchain/blockchain"
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

func (bcn *BlockchainNode) Run() {
	http.HandleFunc("/", bcn.GetChain)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcn.port)), nil))
}
