# Flow Go SDK Demo

## 要件
- Go 1.16
- Flow CLI v0.21.0

## 使い方
### インストール
```sh
go get github.com/onflow/flow-go-sdk@958fc05a220c18276590aa1848907776d0fe24a1
go get github.com/onflow/flow/protobuf/go/flow
go get github.com/golang/protobuf
go get google.golang.org/grpc
```

### エミュレータ起動
```sh
flow emulator
```

### 任意のメッセージへの署名
```sh
go run ./sign-message/main.go
```

