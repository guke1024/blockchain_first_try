package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Add block success!")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("Version: %d\n", block.Version)
		fmt.Printf("Prev block hash: %x\n", block.PrevHash)
		fmt.Printf("Merkel root: %x\n", block.MerkelRoot)
		fmt.Printf("Time stamp: %d\n", block.TimeStamp)
		fmt.Printf("Difficulty: %d\n", block.Difficulty)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Current block hash: %x\n", block.Hash)
		fmt.Printf("Block data: %s\n", block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Printf("==blockchain range end==")
			break
		}
	}
}
