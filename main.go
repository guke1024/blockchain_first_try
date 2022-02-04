package main

func main() {
	bc := NewBlockChain("First")
	cli := CLI{bc}
	cli.Run()
}
