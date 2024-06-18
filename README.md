
# Fastape
Fastape means fast tape. It provides an ultra fast simple data serializer Go module with minimum memory usage.
Fastape just copies data blocks to a byte array and vice versa. Particularly it will shine when size of block is fixed. In golang a data block that contains any combination of numbers, bool, byte and arrays is a fixed sized data block. instead other types including strings, pointers, maps, slices and timers have variable data size. Fastape supports both types simply.

Fastape focuses on maximum throughput and minimum memeory usage. It is specially designed to use in in-memory database, big cache and similar applications that need to serialize/deserialize huge amount of in-memory data. 

## Features
- Cross platform
- ultra fast (first or one of bests based on goserbench benchmark)
- bestest for fixed sized data blocks
- minimum data length overhead
- supports long data size
- supports all data types including basic types, pointers, arrays, maps, slices (except channels and functions)
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
cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.
BenchmarkEncodeVarSize-4                  12097     87.87 ns/op           64 B/op          1 allocs/op
BenchmarkEncodeFixedSize-4                29295     45.54 ns/op           48 B/op          1 allocs/op
BenchmarkEncodeVarSizePreAllocate-4       30975     35.21 ns/op            0 B/op          0 allocs/op
BenchmarkEncodeFixedSizePreAllocate-4     2877007   00.40 ns/op            0 B/op          0 allocs/op
```
## goserbench 
Based on [goserbench](https://github.com/alecthomas/go_serialization_benchmarks) results Enc,Mus-go and Fastape are three winners almost with same results. The winner may be changed by different input data. The following is one of result which shows the efficiency of these compared to fastjson. However both Mus-go and Fastape uses less memory than Enc.
```sh
goos: linux
goarch: amd64
pkg: github.com/alecthomas/go_serialization_benchmarks
cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz
fastape-4   10066           107.1 ns/op              46.00 B/serial        48 B/op          1 allocs/op
benc-4      9090            121.3 ns/op              51.00 B/serial        64 B/op          1 allocs/op
mus-4       8391            133.4 ns/op              46.00 B/serial        48 B/op          1 allocs/op
```
## Compression
Fastape just copies data without any compression, if it's needed it can be used in a next seperate stage after data serialization.

## API 
 Every tape in Fastape provides 3 simple API:

 ```go
 type Tape interface[V any] {
    Sizeof(v V) int
    Roll(v V, bs []byte) (int n, err error)     //Marshall
    Unroll(bs []byte, v *V) (n int, err error)  //Unmarshall     
 }
 ```

## Usage
 Fastape has 6 basic tape and any other new tape can be made based on combinations of these primitives elements:
UnitTape, TimeTape, StringTape, PtrTape, MapTape and SliceTape

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

func (t CompositeTape)Roll(c Composite, bs []byte) (int n, err error){
    k,err = t.stringTape.Roll(c.string,bs[n:]); if err != nil {return}; n+=k
    k,err = t.ptrTape.Roll(c.ptr, bs[n:]); if err != nil {return}; n+=k
    return n, nil
}

func (t compositeTape) Unroll(bs []byte, c *Composite) (n int, err error) {
	  if c == nil { panic("target pointer is nil") }
      var k:=0
	  k,err = t.stringTape.Unroll(bs[n:], &c.string); if err != nil {return}; n+=k
      k,err = t.ptrTape.Unroll(bs[n:], &c.ptr); if err != nil {return}; n+=k
      return n, nil
}

func main() {
    var c,d Composite{}
    c.ptr = &Block{ float64: 3.14 }

    var t CompositeTape
    bs:=make([]byte, t.Sizeof())
    n,err := t.Roll(c, bs)
    print(n, err)
    n,err = t.Unroll(bs, &d)
    print(n, err, d.ptr.float64)
}
```

## license
  **MIT**
