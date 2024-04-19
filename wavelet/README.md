profile

```
go test -cpuprofile cpu.prof -memprofile mem.prof -bench BenchmarkWavelet1MB -run 'Bench.*'
go tool pprof -http=":8888" cpu.prof (or mem.prof)
```
