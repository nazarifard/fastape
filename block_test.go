package fastape

import (
	"testing"

	"github.com/sanity-io/litter"
)

type Block struct {
	Age    int
	Status bool
	Grade  float64
	//comment [16]byte
	Flag byte
	Seq  uint64
}

type BlockTape = UnitTape[Block]

var block = Block{
	Age:    43,
	Status: true,
	Grade:  17.897,
	//comment: [16]byte{'m', 'y', ' ', 't', 'e', 's', 't', ' ', 'n', 'o', 't', 'e'},
	Flag: 'T',
	Seq:  0,
}

func TestBlock(t *testing.T) {
	var blockTape BlockTape
	buff := make([]byte, 1000)
	_, _ = blockTape.Roll(block, buff)

	var newBlock Block
	_, _ = blockTape.Unroll(buff, &newBlock)

	if newBlock != block {
		litter.Dump(block)
		litter.Dump(newBlock)
		t.Errorf("unitTape failed")
	}
}
