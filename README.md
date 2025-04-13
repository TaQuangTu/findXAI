# FindXAI

Find everything at realtime.


## Docker build

```bash
docker buildx build -t findxai:latest .
```

## Development Setup

* Install precommit:

```bash
bash scripts/install_precommit.sh
```

Make sure you see this message as a confirm for a success installation.

```text
✅ Pre-commit setup completed successfully!
```

* gRPC gen Go code:

```bash
protoc --go_out=. --go-grpc_out=. api/search.proto
```

If getting error due to lacking build dependencies. Install Go gRPC generation utilities first:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

And add go bin to PATH
```bash
export PATH=$PATH:$HOME/go/bin
```

Then retry the proto gen command.

* gRPC gen Python code:

```bash
python -m grpc_tools.protoc -I./api --python_out=. --pyi_out=. --grpc_python_out=. content.proto
python -m grpc_tools.protoc -I./api --python_out=. --pyi_out=. --grpc_python_out=. search.proto
```
