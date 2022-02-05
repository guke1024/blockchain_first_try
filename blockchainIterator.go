package main

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

type BlockChainIterator struct {
	db                 *bolt.DB
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		bc.db,
		bc.tail, // It initially points to the last block and changes with the next() function call
	}
}

// Next Iterator belong to blockchain, Next() belong to Iterator
func (it *BlockChainIterator) Next() *Block {
	var block Block
	_ = it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("(iterator)Bucket is nil!")
		}
		blockTmp := bucket.Get(it.currentHashPointer)
		block = Deserialize(blockTmp)
		it.currentHashPointer = block.PrevHash // Move cursor hash left
		return nil
	})
	return &block
}
