# Capy Converter
Converting files, works through grpc

## Run container:
```
docker run -d --restart=unless-stopped --log-opt max-size=5m \
 -p 127.0.0.1:3003:3003 \
 -e MAX_FILE_SIZE_MB="100" \
 --name=capy-converter krol44/capy-converter:latest
```

## Api, proto model:
```
service Converter {
  // converting gif to webm (vp9)
  rpc GifToWebM (GifToWebMType) returns (GifToWebMType) {}
}

message GifToWebMType {
  bytes file = 1;
}
```

## How to use in golang:
1. Run container
2. Import "proto model" into your project
```
import "github.com/krol44/capy-converter/pkg"
```
3. View examples in the **converter_test.go**

### Info for dev:
```
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/model.proto
```