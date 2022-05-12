
### 1 生成pb.go

```
protoc --gofast_out=plugins=grpc:./src ./proto/hello.proto
```

### 2 编译

```
go build client.go hello.pb.go
go build server.go hello.pb.go

```