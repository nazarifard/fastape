package tape

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
		n, _ = c.Marshal(block, bs)
	}
	if n != int(unsafe.Sizeof(block)) {
		b.Errorf("BenchmarkEncodeFixedSizePreAllocate failed!")
	}
}

func BenchmarkEncodeFixedSize(b *testing.B) {
	var c BlockTape
	var n int
	for i := 0; i < b.N; i++ {
		bs := make([]byte, unsafe.Sizeof(*(*Block)(nil)))
		n, _ = c.Marshal(block, bs)
	}
	if n != int(unsafe.Sizeof(block)) {
		b.Errorf("BenchmarkEncodeFixedSize failed!")
	}
}
func BenchmarkFixedSizeDecode(b *testing.B) {
	var v Block
	bs := make([]byte, unsafe.Sizeof(*(*Block)(nil)))
	var c BlockTape
	c.Marshal(block, bs)
	for i := 0; i < b.N; i++ {
		_, _ = c.Unmarshal(bs, &v)
	}
	if v != block {
		b.Errorf("BenchmarkFixedSizeDecode failed!")
	}
}

func BenchmarkEncodeVarSizePreAllocate(b *testing.B) {
	var c compositeTape
	bs := make([]byte, c.Sizeof(compo))
	for i := 0; i < b.N; i++ {
		c.Marshal(compo, bs)
	}
}

func BenchmarkEncodeVarSize(b *testing.B) {
	var c compositeTape
	for i := 0; i < b.N; i++ {
		bs := make([]byte, c.Sizeof(compo))
		c.Marshal(compo, bs)
	}
}
func BenchmarkVarSizeDecode(b *testing.B) {
	bs := make([]byte, 1000)
	var c compositeTape
	c.Marshal(compo, bs)
	var s Composite
	for i := 0; i < b.N; i++ {
		_, _ = c.Unmarshal(bs, &s)
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
	n, _ := bt.Marshal(block, buff)
	_ = ct.Marshal(compo, buff[n:])

	var newBlock Block
	var newCompo Composite
	n, _ = bt.Unmarshal(buff, &newBlock)
	_, _ = ct.Unmarshal(buff[n:], &newCompo)

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
