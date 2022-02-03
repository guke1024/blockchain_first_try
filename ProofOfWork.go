package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000" // appoint difficulty
	tmpInt := big.Int{}
	tmpInt.SetString(targetStr, 16)
	pow.target = &tmpInt
	return &pow
}

func (pow *ProofOfWork) Run() ([]byte, uint64) {
	var nonce uint64
	var hash [32]byte
	block := pow.block
	for {
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkelRoot,
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficulty),
			Uint64ToByte(nonce),
			block.Data,
		}
		blockInfo := bytes.Join(tmp, []byte{})
		hash = sha256.Sum256(blockInfo)
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])          // let hash to big.int
		if tmpInt.Cmp(pow.target) == -1 { // compare generative hash and target
			fmt.Printf("Mining success! hash: %x, nonce: %d\n", hash, nonce)
			return hash[:], nonce
		} else {
			nonce++
		}
	}
}
