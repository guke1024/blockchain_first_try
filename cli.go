package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA "add data to blockchain"
	printChain           "print all blockchain data"
`

// Run go build example => ./example.exe command
func (cli *CLI) Run() {
	args := os.Args // get command
	if len(args) < 2 {
		fmt.Println("Too few parameters")
		fmt.Printf(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Println("Add block:")
		if len(args) == 4 && args[2] == "--data" {
			data := args[3] // get command line data
			cli.AddBlock(data)
		} else {
			fmt.Println("addBlock parameters error")
			fmt.Printf(Usage)
		}
	case "printChain":
		fmt.Println("Print block:")
		cli.PrintBlockChain()
	default:
		fmt.Println("Invalid command")
		fmt.Printf(Usage)
	}
}
