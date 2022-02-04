package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(data) //TODO
	fmt.Println("Add block success!")
}

func (cli *CLI) PrintBlockChain() {
	cli.bc.PrintChain()
	fmt.Println("print blockchain success")
}

func (cli *CLI) PrintBlockChainReverse() {
	bc := cli.bc
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Println("===========================")
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("Prev block hash: %x\n", block.PrevHash)
		fmt.Printf("Merkel root: %x\n", block.MerkelRoot)
		fmt.Printf("Time stamp: %d\n", block.TimeStamp)
		fmt.Printf("Difficulty: %d\n", block.Difficulty)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Current block hash: %x\n", block.Hash)
		fmt.Printf("Block data: %s\n", block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Printf("==blockchain range end==")
			break
		}
	}
}
