package main

import (
	"fmt"
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
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		if len(block.PrevHash) == 0 {
			fmt.Printf("==Blockchain range end==")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {
	if !IsVailAddress(address) {
		fmt.Printf("Address invalid: %s\n", address)
		return
	}
	pubKeyHash := GetPubKeyFromAddress(address)
	utxos := cli.bc.FindUTXOs(pubKeyHash)
	total := 0.0
	for _, utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("%s have balance: %f\n", address, total)
}

func (cli *CLI) Transfer(from, to string, amount float64, miner, data string) {
	if !IsVailAddress(from) {
		fmt.Printf("From address invalid: %s\n", from)
		return
	}
	if !IsVailAddress(to) {
		fmt.Printf("To address invalid: %s\n", to)
		return
	}
	if !IsVailAddress(miner) {
		fmt.Printf("Miner address invalid: %s\n", miner)
		return
	}
	mining := NewMiningTX(miner, data)
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}
	cli.bc.AddBlock([]*Transaction{mining, tx})
	fmt.Println("Transfer success!")
}

func (cli *CLI) CliNewWallet() {
	ws := NewWallets()
	address := ws.CreateWallet()
	fmt.Printf("Address: %s\n", address)

}

func (cli *CLI) ListAddresses() {
	ws := NewWallets()
	addresses := ws.ListAllAddresses()
	for _, address := range addresses {
		fmt.Printf("Address: %s\n", address)
	}
}
