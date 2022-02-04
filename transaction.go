package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 50

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

// IsMining judge whether the current transaction is mining
func (tx *Transaction) IsMining() bool {
	if len(tx.TXInputs) == 1 && len(tx.TXInputs[0].TXid) == 0 && tx.TXInputs[0].Index == -1 {
		return true
	}
	return false
}

// NewMiningTX create a transaction. Mine transaction characteristic: transaction ID and index is not required
func NewMiningTX(address, data string) *Transaction {
	input := TXInput{[]byte{}, -1, data} // Miners don't need to specify sig when mining, sig is usually the name of the ore pool
	output := TXOutput{reward, address}
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()
	return &tx
}

// NewTransaction create a common transaction
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	utxos, resValue := bc.FindNeedUTXOs(from, amount)
	if resValue < amount {
		fmt.Println("Balance is not enough, transaction fail.")
		return nil
	}
	var input []TXInput
	var output []TXOutput
	for id, indexArray := range utxos {
		for _, i := range indexArray { // create transaction input
			input = append(input, TXInput{[]byte(id), int64(i), from})
		}

	}
	output = append(output, TXOutput{amount, to}) // create transaction output

	if resValue > amount { // give change
		output = append(output, TXOutput{resValue - amount, from})
	}
	tx := Transaction{[]byte{}, input, output}
	tx.SetHash()
	return &tx
}
