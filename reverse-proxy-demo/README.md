# 基于 basic auth 的反向代理

本地先起一个反向代理的服务器

```
go run main.go -adress=http://localhost:8090 -username=admin -password=admin7 -port=9002
```

```
go tool pprof --http=:9001 http://localhost:9002/debug/pprof/heap
```