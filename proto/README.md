# Proto (gRPC + gRPC-Gateway)

Generate Go stubs and the gRPC-Gateway from the `.proto` files:

```bash
# From calculation_api/
make proto
```

**Prerequisites:**

1. [buf](https://buf.build/docs/installation) installed and on `PATH`.
2. Protoc plugins on `PATH` (install with Go):
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
   ```
   Ensure `$GOPATH/bin` or `$GOBIN` is in your `PATH`.

After `make proto`, generated files (e.g. `*.pb.go`, `*_grpc.pb.go`, `*.gw.pb.go`) appear next to the source `.proto` files.
