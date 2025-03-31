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

* gRPC gen Python code:

```bash
python -m grpc_tools.protoc -I./api --python_out=. --pyi_out=. --grpc_python_out=. search.proto
```
