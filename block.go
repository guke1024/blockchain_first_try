package main

import (
	"bytes"
	"crypto/sha256"
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

	block.SetHash()
	return &block
}

// SetHash generate hash
func (block *Block) SetHash() {
	tmp := [][]byte{
		Uint64ToByte(block.Version),
		block.PrevHash,
		block.MerkelRoot,
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficulty),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}
