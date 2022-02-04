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
	_ = bc.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("blockBucket"))
		_ = b.ForEach(func(k, v []byte) error {
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

// IsMining judge whether the current transaction is mining
func (tx *Transaction) IsMining() bool {
	if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		if !bytes.Equal(input.TXid, []byte{}) || input.Index != -1 {
			return false
		}
	}
	return true
}

// FindUTXOs find designated address add utxo
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var utxo []TXOutput
	spentOutputs := make(map[string][]int64)
	it := bc.NewIterator() // create block iterator
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			fmt.Printf("current txId: %x\n", tx.TXID)
		JumpOutput:
			for i, output := range tx.TXOutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j { // current ready add output has consumed
							continue JumpOutput // jump
						}
					}
				}
				if output.PubKeyHash == address {
					utxo = append(utxo, output)
				}
			}
			if !tx.IsMining() { // if current transaction is mining, skip directly
				for _, input := range tx.TXInputs {
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray, input.Index)
					}
				}
			} else {
				fmt.Println("This is mining transaction, cancel range.")
			}
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("Blockchain range complete!\n")
			break
		}
	}
	return utxo
}
