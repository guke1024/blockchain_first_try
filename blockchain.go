package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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
	HandleErr("NewBlockChain connect sql fail:\n", err1)
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
	coinbase := NewMiningTX(address, "Genesis block")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// AddBlock add block
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Println("Miner find invalid transactions!")
			return
		}
	}
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
			fmt.Printf("Block data: %s\n", block.Transactions[0].TXInputs[0].PubKey)
			return nil
		})
		return nil
	})
}

// FindUTXOs find designated address add utxo
func (bc *BlockChain) FindUTXOs(PubKeyHash []byte) []TXOutput {
	var utxo []TXOutput
	txs := bc.FindUTXOTransactions(PubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(PubKeyHash, output.PubKeyHash) {
				utxo = append(utxo, output)
			}
		}
	}
	return utxo
}

func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]uint64, float64) {
	utxos := make(map[string][]uint64) // find utxos what need
	var calc float64                   // sum utxos
	txs := bc.FindUTXOTransactions(senderPubKeyHash)
	for _, tx := range txs {
		for i, output := range tx.TXOutputs {
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				if calc < amount {
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], uint64(i)) // add utxo
					calc += output.Value                                               // sum current utxo
					if calc >= amount {
						return utxos, calc
					}
				} else {
					fmt.Printf("Transfer amount not satisfied, current total %f", calc)
				}
			}
		}
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransactions(senderPubKeyHash []byte) []*Transaction {
	var txs []*Transaction // store all transaction
	spentOutputs := make(map[string][]int64)
	it := bc.NewIterator() // create block iterator
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			//fmt.Printf("current txId: %x\n", tx.TXID)
		JumpOutput:
			for i, output := range tx.TXOutputs {
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j { // current ready add output has consumed
							continue JumpOutput // jump
						}
					}
				}
				if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
					txs = append(txs, tx) // all transactions involving utxo
				}
			}
			if !tx.IsMining() { // if current transaction is mining, skip directly
				for _, input := range tx.TXInputs {
					if bytes.Equal(HashPubKey(input.PubKey), senderPubKeyHash) {
						spentOutputs[string(input.TXid)] = append(spentOutputs[string(input.TXid)], input.Index)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return txs
}

func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	it := bc.NewIterator()
	for {
		block := it.Next()
		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("invalid transaction")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)
	for _, inp := range tx.TXInputs {
		tx, err1 := bc.FindTransactionByTXid(inp.TXid)
		HandleErr("", err1)
		prevTXs[string(inp.TXid)] = tx
	}
	tx.Sign(privateKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsMining() {
		return true
	}
	prevTXs := make(map[string]Transaction)
	for _, inp := range tx.TXInputs {
		tx, err1 := bc.FindTransactionByTXid(inp.TXid)
		HandleErr("", err1)
		prevTXs[string(inp.TXid)] = tx
	}
	return tx.Verify(prevTXs)
}
