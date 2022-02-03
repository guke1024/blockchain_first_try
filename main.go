package main

func main() {
	bc := NewBlockChain()
	cli := CLI{bc}
	cli.Run()
	//bc.AddBlock("first block")
	//bc.AddBlock("second block")
}
