package main

import (
	"fmt"
	"time"
)

func (cli *CLI) PrintBlockChain() {
	cli.bc.PrintChain()
	fmt.Println("Print blockchain success")
}

func (cli *CLI) PrintBlockChainReverse() {
	bc := cli.bc
	it := bc.NewIterator()
	for {
		block := it.Next()
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Println("===========================")
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("Prev block hash: %x\n", block.PrevHash)
		fmt.Printf("Merkel root: %x\n", block.MerkelRoot)
		fmt.Printf("Time stamp: %s\n", timeFormat)
		fmt.Printf("Difficulty: %d\n", block.Difficulty)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Current block hash: %x\n", block.Hash)
		fmt.Printf("Block data: %s\n", block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Printf("==Blockchain range end==")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s have balance: %f\n", address, total)
}

func (cli *CLI) Transfer(from, to string, amount float64, miner, data string) {
	mining := NewMiningTX(miner, data)
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}
	cli.bc.AddBlock([]*Transaction{mining, tx})
	fmt.Println("Transfer success!")
}

func (cli *CLI) CliNewWallet() {
	wallet := NewWallet()
	address := wallet.NewAddress()
	fmt.Printf("Private Key: %v\n", wallet.Private)
	fmt.Printf("Public Key: %v\n", wallet.Public)
	fmt.Printf("Address: %s\n", address)
}
