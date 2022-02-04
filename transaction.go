package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

const reward = 12.5

type Transaction struct {
	TXID      []byte // transaction id
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXid  []byte
	Index int64 // quote output index
	Sig   string
}

type TXOutput struct {
	Value      float64
	PubKeyHash string // lock script
}

// SetHash set transaction ID
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

// NewCoinbaseTX create a transaction. Mine transaction characteristic: transaction ID and index is not required
func NewCoinbaseTX(address, data string) *Transaction {
	input := TXInput{[]byte{}, -1, data} // Miners don't need to specify sig when mining, sig is usually the name of the ore pool
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
}
