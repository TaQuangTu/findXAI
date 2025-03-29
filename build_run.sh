go build -o server cmd/server/main.go
export $(grep -v '^#' .env | xargs)
./server