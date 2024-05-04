package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/jvsena42/go_blockchain/blockchain"
	"github.com/jvsena42/go_blockchain/utils"
	"github.com/jvsena42/go_blockchain/wallet"
)

const pathToTemplateDir = "templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port: port, gateway: gateway}
}

func (ws *WalletServer) Port() uint16 {
	return ws.port
}

func (ws *WalletServer) Gateway() string {
	return ws.gateway
}

func (ws *WalletServer) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		t, _ := template.ParseFiles(path.Join(pathToTemplateDir, "index.html"))
		t.Execute(w, "")
	default:
		log.Println("/Index Error: Invalid http request", r.Method)
	}
}

func (ws *WalletServer) Wallet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Header().Add("Content-Type", "application/json")
		myWallet := wallet.NewWallet()
		marshalWallet, _ := myWallet.MarshalJSON()
		io.WriteString(w, string(marshalWallet[:]))

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("/wallet Error: Invalid http request", r.Method)
	}
}

func (ws *WalletServer) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var t wallet.TransactionRequest
		err := decoder.Decode(&t)
		if err != nil {
			log.Printf("ERROR: %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if !t.Validate() {
			log.Println("ERROR: missing fields")
			io.WriteString(w, string(utils.JsonStatus("ERROR: missing fields")))
			return
		}

		publicKey := utils.StringToPublicKey(*t.SenderPublicKey)
		privateKey := utils.StringToPrivateKey(*t.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*t.Value, 32)

		if err != nil {
			log.Println("ERROR: parsing value")
			io.WriteString(w, string(utils.JsonStatus("Error: invalid value")))
			return
		}

		value32 := float32(value)

		w.Header().Add("Content-Type", "application/json")

		transaction := wallet.NewTransaction(privateKey, publicKey, *t.SenderBlockchainAddress, *t.RecipientBlockchainAddress, value32)
		signature := transaction.GenerateSignature()
		signatureStr := signature.String()

		bt := &blockchain.TransactionRequest{
			SenderBlockchainAddress:    t.SenderBlockchainAddress,
			RecipientBlockchainAddress: t.RecipientBlockchainAddress,
			SenderPublicKey:            t.SenderPublicKey,
			Value:                      &value32,
			Signature:                  &signatureStr,
		}

		m, err := json.Marshal(bt)

		if err != nil {
			log.Printf("/Trancasctions ERROR: parsing json %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		buff := bytes.NewBuffer(m)
		resp, err := http.Post(ws.Gateway()+"/transactions", "application/json", buff)

		if err != nil {
			log.Printf("/Trancasctions ERROR: request %v", err)
			io.WriteString(w, string(utils.JsonStatus("fail")))
			return
		}

		if resp.StatusCode == 201 {
			io.WriteString(w, string(utils.JsonStatus("success")))
			return
		} else {
			io.WriteString(w, string(utils.JsonStatus("error")))
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		log.Println("/transactions ERROR: Invalid HTTP method", r.Method)
	}
}

func (ws *WalletServer) Run() {
	http.HandleFunc("/", ws.Index)
	http.HandleFunc("/wallet", ws.Wallet)
	http.HandleFunc("/transactions", ws.CreateTransaction)

	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(ws.port)), nil))
}
