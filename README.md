# Echo

A TCP server that will efficiently write back all received bytes.

## Benchmark

`BenchmarkKnownCount` is the fastest benchmark, but there are other, more realistic, implementations available.
```
go run app/main.go -port=1122
go test -bench=BenchmarkKnownCount -benchmem -count=5
```

### Results
Running in batches of 64KB:
```
goos: linux
goarch: amd64
pkg: github.com/daulet/echo
cpu: AMD Ryzen 9 5900HS with Radeon Graphics        
BenchmarkKnownCount-16          1000000000               0.04737 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.06366 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.06375 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.05343 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.06283 ns/op         0 B/op          0 allocs/op
PASS
ok      github.com/daulet/echo  0.418s
```

### BPF-Echo
[eBPF echo server](https://github.com/path-network/bpf-echo) implementation doesn't seem to be able to handle batch size above 16KB, but it seems to be a bit more efficient:
```
goos: linux
goarch: amd64
pkg: github.com/daulet/echo
cpu: AMD Ryzen 9 5900HS with Radeon Graphics        
BenchmarkKnownCount-16          1000000000               0.03328 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.02866 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.02907 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.03115 ns/op         0 B/op          0 allocs/op
BenchmarkKnownCount-16          1000000000               0.03171 ns/op         0 B/op          0 allocs/op
PASS
ok      github.com/daulet/echo  0.224s
```