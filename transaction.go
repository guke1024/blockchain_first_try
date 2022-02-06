package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
)

const reward = 50

type Transaction struct {
	TXID      []byte // transaction id
	TXInputs  []TXInput
	TXOutputs []TXOutput
}

type TXInput struct {
	TXid      []byte
	Index     int64  // quote output index
	Signature []byte //
	PubKey    []byte
}

type TXOutput struct {
	Value      float64
	PubKeyHash []byte // lock script
}

func (output *TXOutput) Lock(address string) {
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{Value: value}
	output.Lock(address)
	return &output
}

// SetHash set transaction ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err1 := encoder.Encode(tx)
	HandleErr("SetHash Encode!", err1)
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
	input := TXInput{[]byte{}, -1, nil, []byte(data)} // Miners don't need to specify sig when mining, sig is usually the name of the ore pool
	output := NewTXOutput(reward, address)
	tx := Transaction{[]byte{}, []TXInput{input}, []TXOutput{*output}}
	tx.SetHash()
	return &tx
}

// NewTransaction create a common transaction
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	ws := NewWallets()
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Println("No wallet found for this address, create transaction fail!")
		return nil
	}
	publicKey := wallet.Public
	privateKey := wallet.Private

	utxos, resValue := bc.FindNeedUTXOs(HashPubKey(publicKey), amount)
	if resValue < amount {
		fmt.Println("Balance is not enough, transaction fail.")
		return nil
	}
	var input []TXInput
	var output []TXOutput
	for id, indexArray := range utxos {
		for _, i := range indexArray { // create transaction input
			input = append(input, TXInput{[]byte(id), int64(i), nil, publicKey})
		}

	}
	output = append(output, *NewTXOutput(amount, to)) // create transaction output

	if resValue > amount { // give change
		output = append(output, *NewTXOutput(resValue-amount, from))
	}
	tx := Transaction{[]byte{}, input, output}
	tx.SetHash()
	bc.SignTransaction(&tx, privateKey)
	return &tx
}

func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsMining() {
		return
	}
	txCopy := tx.TrimmedCopy()
	for i, input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("Quote transaction invalid!")
		}
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		r, s, err1 := ecdsa.Sign(rand.Reader, privateKey, signDataHash)
		HandleErr("Sign ecdsa.Sign:\n", err1)
		signature := append(r.Bytes(), s.Bytes()...)
		tx.TXInputs[i].Signature = signature
	}
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput
	for _, input := range tx.TXInputs {
		inputs = append(inputs, TXInput{input.TXid, input.Index, nil, nil})
	}
	for _, output := range tx.TXOutputs {
		outputs = append(outputs, output)
	}
	return Transaction{tx.TXID, inputs, outputs}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsMining() {
		return true
	}
	txCopy := tx.TrimmedCopy()
	for i, input := range tx.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("Quote transaction invalid!")
		}
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		dataHash := txCopy.TXID
		signature := input.Signature
		pubKey := input.PubKey
		r := big.Int{}
		s := big.Int{}
		r.SetBytes(signature[:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])
		X := big.Int{}
		Y := big.Int{}
		X.SetBytes(pubKey[:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrigin := ecdsa.PublicKey{Curve: elliptic.P256(), X: &X, Y: &Y}
		if !ecdsa.Verify(&pubKeyOrigin, dataHash, &r, &s) {
			return false
		}
	}
	return true
}
