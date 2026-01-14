# Benchmarks

Numbers depend on hardware and workload. The following results are a snapshot from the original author.

## Local Go benchmark

```
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

Based on https://github.com/alecthomas/go_serialization_benchmarks.

```
goos: linux
goarch: amd64
pkg: github.com/alecthomas/go_serialization_benchmarks
cpu: Intel(R) Core(TM) i7-3537U CPU @ 2.00GHz

#  name                          iter     ns/op   size   B/op  allocs/op
0  fastape_reuse-4               4624285  51.93   46.00  0     0
1  baseline/unsafe_reuse-4       4048184  61.19   47.00  0     0
2  mus/unsafe_reuse-4            3476158  68.79   49.00  0     0
3  gencode/unsafe_reuse-4        3286796  72.64   46.00  0     0
4  wellquite/bebop/reuse-4       2925176  82.63   55.00  0     0
5  200sc/bebop/reuse-4           3159244  84.50   55.00  0     0
6  fastape-4                     2091486  118.8   46.00  48    1
7  baseline-4                    1746391  121.6   47.00  48    1
8  benc/usafe-4                  1746405  135.8   51.00  64    1
9  benc-4                        1758759  136.9   51.00  64    1
10 200sc/bebop-4                 1547842  151.7   55.00  64    1
11 mus-4                         1559601  152.6   46.00  48    1
12 wellquite/bebop-4             1532216  159.2   55.00  64    1
13 msgp-4                        1249694  192.3   97.00  128   1
14 colfer-4                      1000000  221.0   51.09  64    1
15 gencode-4                     957394   232.6   53.00  80    2
16 calmh/xdr-4                   1000000  250.5   60.00  64    1
17 gogo/protobuf-4               1291285  284.1   53.00  64    1
18 shamaton/msgpackgen/array-4   690806   377.3   50.00  144   2
19 shamaton/msgpackgen/map-4     510366   392.3   92.00  176   2
20 gotiny-4                      458545   509.8   48.00  168   5
21 hprose2-4                     441703   522.0   85.27  0     0
22 dedis/protobuf-4              324448   746.3   52.00  144   7
23 ikea-4                        315805   749.9   55.00  72    8
24 pulsar-4                      261908   877.6   51.62  304   7
25 jsoniter-4                    225681   930.9   141.4  200   3
26 shamaton/msgpack/array-4      248595   942.3   50.00  160   4
27 flatbuffers-4                 212164   947.1   95.21  376   10
28 msgpack-4                     233332   957.3   92.00  264   4
29 shamaton/msgpack/map-4        225324   1045    92.00  192   4
30 hprose-4                      217099   1098    85.26  453   8
31 bson-4                        184314   1100    110.0  376   10
32 avro2/binary-4                172772   1275    47.00  464   9
33 ugorji/msgpack-4              171700   1307    91.00  1240  3
34 easyjson-4                    164289   1361    150.8  976   7
35 ugorji/binc-4                 156680   1459    95.00  1256  4
36 davecgh/xdr-4                 144692   1499    92.00  392   20
37 mongobson-4                   133244   1641    110.0  240   9
38 json-4                        129004   1687    151.8  208   2
39 alecthomas/binary-4           116653   1884    61.00  360   27
40 fastjson/reuse-4              103352   2227    133.8  1360  7
41 goavro-4                      87268    2672    47.00  584   18
42 fastjson-4                    74778    2904    133.8  1864  13
43 capnproto-4                   78266    2963    96.00  4392  6
44 sereal-4                      69087    3206    142.0  1104  22
45 avro2/text-4                  51373    4409    133.8  1320  20
46 ssz-4                         46316    4715    55.00  416   45
47 gob-4                         37384    6206    172.6  1744  37
48 gogo/jsonpb-4                 13254    17768   125.5  2747  80
```
