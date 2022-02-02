package main

import "crypto/sha256"

type Block struct {
	PrevHash []byte
	Hash     []byte
	Data     []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		PrevHash: prevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

// SetHash generate hash
func (block *Block) SetHash() {
	blockInfo := append(block.PrevHash, block.Data...)
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
