package main

import (
	"bytes"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
)

type BlockChain struct {
	db   *bolt.DB
	tail []byte // store last block hash
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

// NewBlockChain create block and add to blockchain
func NewBlockChain(address string) *BlockChain {
	var lastHash []byte
	db, err1 := bolt.Open(blockChainDb, 0600, nil)
	if err1 != nil {
		log.Panic("connect sql fail")
	}
	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			bucket, _ = tx.CreateBucket([]byte(blockBucket))
			genesisBlock := GenesisBlock(address)
			_ = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			_ = bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db, lastHash}
}

// GenesisBlock def Genesis block
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, "Genesis block")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// AddBlock add block
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// get last block and it's hash
	db := bc.db
	lastHash := bc.tail
	_ = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket is nil")
		}
		block := NewBlock(txs, lastHash)
		_ = bucket.Put(block.Hash, block.Serialize())
		_ = bucket.Put([]byte("LastHashKey"), block.Hash)
		bc.tail = block.Hash
		return nil
	})
}

func (bc *BlockChain) PrintChain() {
	blockHeight := 0
	bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))
		b.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte("LastHashKey")) {
				return nil
			}
			block := Deserialize(v)
			fmt.Printf("=============== block height: %d ===============\n", blockHeight)
			blockHeight++
			fmt.Printf("Version: %d\n", block.Version)
			fmt.Printf("Prev block hash: %x\n", block.PrevHash)
			fmt.Printf("Merkel root: %x\n", block.MerkelRoot)
			fmt.Printf("Time stamp: %d\n", block.TimeStamp)
			fmt.Printf("Difficulty: %d\n", block.Difficulty)
			fmt.Printf("Nonce: %d\n", block.Nonce)
			fmt.Printf("Current block hash: %x\n", block.Hash)
			fmt.Printf("Block data: %s\n", block.Transactions[0].TXInputs[0].Sig)
			return nil
		})
		return nil
	})
}
