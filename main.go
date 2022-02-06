package main

func main() {
	bc := NewBlockChain("1JLyj9EK6sEv2NpVFV74NvEmgakc2uqmUj")
	cli := CLI{bc}
	cli.Run()
}
