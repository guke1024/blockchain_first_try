package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

// Block def block struct
type Block struct {
	Version    uint64
	PrevHash   []byte
	MerkelRoot []byte
	TimeStamp  uint64
	Difficulty uint64
	Nonce      uint64

	// following data is not in real blockchain
	Hash []byte
	Data []byte
}

func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err1 := binary.Write(&buffer, binary.BigEndian, num)
	if err1 != nil {
		log.Panic(err1)
	}
	return buffer.Bytes()
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}

	//block.SetHash()
	pow := NewProofOfWork(&block)
	// select nonce, keep hashing
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	return &block
}
