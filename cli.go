package main

import (
	"fmt"
	"os"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA           "add data to blockchain"
	printChain                     "forward print all blockchain data"
	printChainR                    "Reverse print all blockchain data"
	getBalance --address ADDRESS   "obtain designated address balance"
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
		fmt.Println("forward print all blockchain data:")
		cli.PrintBlockChain()
	case "printChainR":
		fmt.Println("Reverse print all blockchain data:")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Println("obtain designated address balance:")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3] // get command line data
			cli.GetBalance(address)
		} else {
			fmt.Println("getBalance parameters error")
			fmt.Printf(Usage)
		}
	default:
		fmt.Println("Invalid command")
		fmt.Printf(Usage)
	}
}
