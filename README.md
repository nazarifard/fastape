
# Fastape
Fastape means fast tape. It provides an ultra fast simple data serializer Go module.
Fastape just copies data blocks to a byte array and vice versa. Particularly it will shine when size of block is fixed. In golang a data block that contains any combination of numbers, bool, byte and arrays is a fixed sized data block. instead other types including strings, pointers, maps, slices and timers have variable data size. Fastape supports both types simply.

Fastape just uses golang "copy" api and doesn't use any other specific data serializer algorithms,less complexity means less errors.

## Features
- Cross platform
- ultra fast (one of bests based on goserbench)
- bestest for fixed sized data blocks
- minimum data length overhead
- supports long data size
- supports basic types, pointers, arrays, maps, slices
- easy to use with just 6 generic APIs

## Installation
```
go get github.com/nazarifard/fastape
```

## Benchmark
```sh
# go test -bench=. -benchmem 
goos: linux
goarch: amd64
pkg: fastape
cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
BenchmarkEncodeFixedSizePreAllocate-4  723041300     1.471 ns/op           0 B/op          0 allocs/op
BenchmarkEncodeFixedSize-4             851389168     1.453 ns/op           0 B/op          0 allocs/op
BenchmarkFixedSizeDecode-4             650878077     1.665 ns/op           0 B/op          0 allocs/op
BenchmarkEncodeVarSizePreAllocate-4    20057034      60.84 ns/op           0 B/op          0 allocs/op
BenchmarkEncodeVarSize-4               8013882       146.5 ns/op           64 B/op         1 allocs/op
BenchmarkVarSizeDecode-4               24088782      43.47 ns/op           0 B/op          0 allocs/op
```
## goserbench 
Based on [goserbench](https://github.com/alecthomas/go_serialization_benchmarks) results Enc,Mus-go and Fastape are three winners by a small margin. The winner may be changed by different input data, in most cases Enc is winner yet. Howere input data isn`t random realy therefore it isn't a fair play. for example Benc simply converts int(Siblings) to uint16 or uses two bytes to store size of strings while string size is unlimitted in Golang. The following is one of result which shows the efficiency of these compared to fastjson.
```sh
goos: linux
goarch: amd64
pkg: github.com/alecthomas/go_serialization_benchmarks
cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
marshal/fastjson-4          591217   2287.0  ns/op   1360 B/op    7 allocs/op    133.7 B/serial
marshal/benc/usafe-4       9664924    115.7  ns/op     64 B/op    1 allocs/op    51.00 B/serial
marshal/mus/unsafe-4       9648376    122.7  ns/op     64 B/op    1 allocs/op    49.00 B/serial
marshal/fastape-4             7793353    153.1  ns/op     64 B/op    1 allocs/op    55.00 B/serial
                                                                
unmarshal/fastjson-4        432190   2789.00 ns/op   1800 B/op   11 allocs/op
unmarshal/benc/usafe-4    23514747     51.97 ns/op   0    B/op    0 allocs/op
unmarshal/mus/unsafe-4    15785749     73.65 ns/op   0    B/op    0 allocs/op
unmarshal/fastape-4          16388287     72.01 ns/op   0    B/op    0 allocs/op
```



## Compression
Fastape just copies data without any compression, if it's needed it can be used in a next seperate stage after data serialization.
## Usage
Fastape has 6 basic tape and any other new tape can be made based on combinations of these primitives elements:
UnitTape, TimeTape, StringTape, PtrTape, MapTape, SliceTape

examle:
```go
package main
import "github.com/nazarifard/fastape"

type Block struct {
	bool
	byte
	int
	float64
	fiveF [5]float32
	ten   [10]byte
}
type Composite struct {
    string 
    ptr *Block
    //time time.Time
    //map map[int]Block
    //slice []Block
}

type BlockTape=fastape.Unit[Block]
type CompositeTape struct {
    stringTape fastape.StringTape
    ptrTape fastape.PtrTape[Block]
    //timeTape fastape.TimeTape
    //map fastape.MapTape[int,Block, Unit[int], BlockTape]
    //Slice fastape.Slice[Block]
}
func (t CompositeTape)Sizeof(c Composite) int {
    return t.stringTape.Sizeof() + 
           t.ptrTape.Sizeof() 
}
func (t CompositeTape)Marshal(c Composite, bs []byte) (int n, err error){
    k,err = t.stringTape.Marshal(c.string,bs[n:]); if err != nil {return}; n+=k
    k,err = t.ptrTape.Marshal(c.ptr, bs[n:]); if err != nil {return}; n+=k
    return n, nil
}
func (t compositeTape) Unmarshal(bs []byte, c *Composite) (n int, err error) {
	  if c == nil { panic("target pointer is nil") }
      var k:=0
	  k,err = t.stringTape.Marshal(bs[n:], &c.string); if err != nil {return}; n+=k
      k,err = t.ptrTape.Marshal(bs[n:], &c.ptr); if err != nil {return}; n+=k
      return n, nil
}

func main() {
    var c,d Composite{}
    c.ptr = &Block{ float64: 3.14 }

    var t CompositeTape
    bs:=make([]byte, t.Sizeof())
    n,err := t.Marshal(c, bs)
    print(n, err)
    n,err = t.Unmarshal(bs, &d)
    print(n, err, d.ptr.float64)
}
```

## license
  **MIT**
