package main

import (
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain                         "Forward print all blockchain data"
	printChainR                        "Reverse print all blockchain data"
	getBalance --address ADDRESS       "Obtain designated address balance"
	transfer FROM TO AMOUNT MINER DATA "FROM transfers AMOUNT to TO, MINER mine and write to data"
	newWallet                          "Create a new wallet"
	listAddresses                      "List all addresses"
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
	case "printChain":
		fmt.Println("Forward print all blockchain data:")
		cli.PrintBlockChain()
	case "printChainR":
		fmt.Println("Reverse print all blockchain data:")
		cli.PrintBlockChainReverse()
	case "getBalance":
		fmt.Println("Obtain designated address balance:")
		if len(args) == 4 && args[2] == "--address" {
			address := args[3] // get command line data
			cli.GetBalance(address)
		} else {
			fmt.Println("GetBalance parameters error!")
			fmt.Printf(Usage)
		}
	case "transfer":
		fmt.Println("Begin transfer...")
		if len(args) != 7 {
			fmt.Println("Transfer parameters error!")
			fmt.Printf(Usage)
			return
		}
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64) // string change to float64
		miner := args[5]
		data := args[6]
		cli.Transfer(from, to, amount, miner, data)
	case "newWallet":
		fmt.Println("Create a new wallet:")
		cli.CliNewWallet()
	case "listAddresses":
		fmt.Println("List all addresses:")
		cli.ListAddresses()
	default:
		fmt.Println("Invalid command")
		fmt.Printf(Usage)
	}
}
