package main

func main() {
	bc := NewBlockChain("Creator")
	cli := CLI{bc}
	cli.Run()
}
