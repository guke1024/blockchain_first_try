package main

type BlockChain struct {
	blocks []*Block
}

// NewBlockChain create block and add to blockchain
func NewBlockChain() *BlockChain {
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

// GenesisBlock def Genesis block
func GenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}

// AddBlock add block
func (bc *BlockChain) AddBlock(data string) {
	// get last block and it's hash
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevHash := lastBlock.Hash
	block := NewBlock(data, prevHash)
	bc.blocks = append(bc.blocks, block)
}
