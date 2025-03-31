package main

import (
	"findx/config"
	"findx/internal/lockdb"
	"findx/internal/server"
	"findx/internal/system"
	"findx/pkg/protogen"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	rootClosers = make([]func() error, 0)
)

func RegisterRootCloser(closer func() error) {
	rootClosers = append(rootClosers, closer)
}

func main() {
	cfg := config.NewConfig()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lockDb, err := lockdb.NewLockDbRedis(cfg.REDIS_LOCKDB_DNS)
	if err != nil {
		return
	}
	rateLimiter, err := lockdb.NewOurRateLimit(cfg.REDIS_LOCKDB_DNS)
	if err != nil {
		return
	}
	s := grpc.NewServer()
	searchServer := server.NewSearchServer(cfg, lockDb, rateLimiter)
	protogen.RegisterSearchServiceServer(s, searchServer)

	defer system.SafeClose()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
