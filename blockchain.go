package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

type BlockChain struct {
	db   bolt.DB
	tail []byte // store last block hash
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

// NewBlockChain create block and add to blockchain
func NewBlockChain() *BlockChain {
	var lastHash []byte
	db, err1 := bolt.Open(blockChainDb, 0600, nil)
	if err1 != nil {
		log.Panic("connect sql fail")
	}
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(blockBucket))
			genesisBlock := GenesisBlock()
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{*db, lastHash}
}

// GenesisBlock def Genesis block
func GenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

// AddBlock add block
func (bc *BlockChain) AddBlock(data string) {
	// get last block and it's hash
	db := bc.db
	lastHash := bc.tail
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket is nil")
		}
		block := NewBlock(data, lastHash)
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)
		bc.tail = block.Hash
		return nil
	})
}
