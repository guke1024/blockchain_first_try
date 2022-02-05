package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.LoadFile()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	ws.WalletsMap[address] = wallet
	ws.SaveToFile()
	return address
}

func (ws *Wallets) SaveToFile() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err1 := encoder.Encode(ws)
	HandleErr("SaveToFile Encode!", err1)
	err2 := ioutil.WriteFile(walletFile, buffer.Bytes(), 0600)
	HandleErr("SaveToFile WriteFile!", err2)
}

func (ws *Wallets) LoadFile() {
	_, err1 := os.Stat(walletFile)
	if os.IsNotExist(err1) {
		return
	}
	content, err2 := ioutil.ReadFile(walletFile)
	HandleErr("LoadFile ioutil.ReadFile:\n", err2)
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(content))
	var wsLocal Wallets
	err3 := decoder.Decode(&wsLocal)
	HandleErr("LoadFile decoder.Decode:\n", err3)
	ws.WalletsMap = wsLocal.WalletsMap
}

func (ws *Wallets) ListAllAddresses() []string {
	var addresses []string
	for address := range ws.WalletsMap {
		addresses = append(addresses, address)
	}
	return addresses
}
