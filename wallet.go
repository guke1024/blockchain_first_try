package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	Public  []byte
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err1 := ecdsa.GenerateKey(curve, rand.Reader)
	HandleErr("NewWallet ecdsa.GenerateKey:\n", err1)
	publicKeyOrig := privateKey.PublicKey
	publicKey := append(publicKeyOrig.X.Bytes(), publicKeyOrig.Y.Bytes()...)
	return &Wallet{privateKey, publicKey}
}

func HashPubKey(data []byte) []byte { // RIPEMD160(sha256(public key))
	hash := sha256.Sum256(data)
	rip160Hash := ripemd160.New() // encoder
	_, err1 := rip160Hash.Write(hash[:])
	HandleErr("HashPubKey ripemd160.Write:\n", err1)
	rip160HashValue := rip160Hash.Sum(nil) // rip160HashValue = public key hash (20bytes data)
	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	hash1 := sha256.Sum256(data)     // sha256(public key hash
	hash2 := sha256.Sum256(hash1[:]) // sha256(sha256(public key hash)
	checkCode := hash2[:4]           // take the first four bytes as the check code
	return checkCode

}

func (w *Wallet) NewAddress() string {
	publicKey := w.Public
	rip160HashValue := HashPubKey(publicKey)
	version := byte(00)
	payload := append([]byte{version}, rip160HashValue...) // 21bytes data
	checkCode := CheckSum(payload)
	payload = append(payload, checkCode...) // 25bytes data
	address := base58.Encode(payload)
	return address
}
