package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Transaction struct {
	TXID     []byte // transaction id
	TXInputs []TXInput
	TXOutput []TXOutput
}

type TXInput struct {
	TXid  []byte
	Index int64
	Sig   string
}

type TXOutput struct {
	value      float64
	PukKeyHash string // lock script
}

func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err1 := encoder.Encode(tx)
	if err1 != nil {
		log.Panic(err1)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}
