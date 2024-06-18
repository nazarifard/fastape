package fastape

import (
	"testing"
	"unsafe"

	"github.com/sanity-io/litter"
)

func BenchmarkEncodeFixedSizePreAllocate(b *testing.B) {
	bs := make([]byte, unsafe.Sizeof(*(*Block)(nil)))
	var c BlockTape
	var n int
	for i := 0; i < b.N; i++ {
		n, _ = c.Roll(block, bs)
	}
	if n != int(unsafe.Sizeof(block)) {
		b.Errorf("BenchmarkEncodeFixedSizePreAllocate failed!")
	}
}

func BenchmarkEncodeFixedSize(b *testing.B) {
	var c BlockTape
	var n int
	var bs []byte
	for i := 0; i < b.N; i++ {
		bs = make([]byte, unsafe.Sizeof(*(*Block)(nil)))
		n, _ = c.Roll(block, bs)
	}
	if n != int(unsafe.Sizeof(block)) {
		b.Errorf("BenchmarkEncodeFixedSize failed!")
	}
	_ = bs
}
func BenchmarkFixedSizeDecode(b *testing.B) {
	var v Block
	bs := make([]byte, unsafe.Sizeof(*(*Block)(nil)))
	var c BlockTape
	c.Roll(block, bs)
	for i := 0; i < b.N; i++ {
		_, _ = c.Unroll(bs, &v)
	}
	if v != block {
		b.Errorf("BenchmarkFixedSizeDecode failed!")
	}
}

func BenchmarkEncodeVarSizePreAllocate(b *testing.B) {
	var c compositeTape
	bs := make([]byte, c.Sizeof(compo))
	for i := 0; i < b.N; i++ {
		c.Roll(compo, bs)
	}
}

func BenchmarkEncodeVarSize(b *testing.B) {
	var c compositeTape
	for i := 0; i < b.N; i++ {
		bs := make([]byte, c.Sizeof(compo))
		n := c.Roll(compo, bs)
		_ = n
	}
}
func BenchmarkVarSizeDecode(b *testing.B) {
	bs := make([]byte, 1000)
	var c compositeTape
	c.Roll(compo, bs)
	var s Composite
	for i := 0; i < b.N; i++ {
		_, _ = c.Unroll(bs, &s)
		// if s.Siblings != sample.Siblings {
		// 	panic("..............................")
		// }
	}
	if s.Siblings != compo.Siblings {
		b.Errorf("BenchmarkDecode failed!")
	}
}

func TestMain(m *testing.M) {
	litter.Config.HidePrivateFields = false
	m.Run()
}

func TestChainBlockComposite(t *testing.T) {
	var bt BlockTape
	var ct compositeTape

	buff := make([]byte, 1000)
	n, _ := bt.Roll(block, buff)
	_ = ct.Roll(compo, buff[n:])

	var newBlock Block
	var newCompo Composite
	n, _ = bt.Unroll(buff, &newBlock)
	_, _ = ct.Unroll(buff[n:], &newCompo)

	if newBlock != block {
		litter.Dump(block)
		litter.Dump(newBlock)

		t.Errorf("TestChainBlockComposite chain1 failed!")
	}
	if newCompo.Name != compo.Name ||
		newCompo.Money != compo.Money ||
		newCompo.Phone != compo.Phone ||
		!newCompo.BirthDay.Equal(compo.BirthDay) ||
		newCompo.Siblings != compo.Siblings ||
		newCompo.Spouse != compo.Spouse {

		litter.Dump(compo)
		litter.Dump(newCompo)
		t.Errorf("TestChainBlockComposite chain2 failed!")
	}
}
