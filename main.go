package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("first block")
	bc.AddBlock("second block")

	it := bc.NewIterator()
	for {
		block := it.Next()
		fmt.Printf("Prev block hash: %x\n", block.PrevHash)
		fmt.Printf("Current block hash: %x\n", block.Hash)
		fmt.Printf("Block data: %s\n", block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Printf("blockchain range end")
			break
		}
	}
}
