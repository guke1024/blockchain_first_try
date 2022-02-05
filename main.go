package main

func main() {
	bc := NewBlockChain("1FeV7iAUYVxcZQ7Xc8j5Pj3WSSUguijTCC")
	cli := CLI{bc}
	cli.Run()
}
