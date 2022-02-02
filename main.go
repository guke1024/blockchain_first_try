package main

import (
	"fmt"
)

func main() {
	bc := NewBlockChain()
	bc.AddBlock("first block")
	bc.AddBlock("second block")
	for i, block := range bc.blocks {
		fmt.Println("Current block height: ", i)
		fmt.Printf("Prev block hash: %x\n", block.PrevHash)
		fmt.Printf("Current block hash: %x\n", block.Hash)
		fmt.Printf("Block data: %s\n", block.Data)
	}

}
