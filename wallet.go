package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	Public  []byte
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err1 := ecdsa.GenerateKey(curve, rand.Reader)
	if err1 != nil {
		log.Panic()
	}
	publicKeyOrig := privateKey.PublicKey
	publicKey := append(publicKeyOrig.X.Bytes(), publicKeyOrig.Y.Bytes()...)
	return &Wallet{privateKey, publicKey}

}
