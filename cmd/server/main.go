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

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.PORT, cfg.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	lockDb, err := lockdb.NewLockDbRedis(cfg.REDIS_LOCKDB_DSN)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	rateLimiter, err := lockdb.NewOurRateLimit(cfg.REDIS_RATE_LIMITDB_DSN)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	searchServer := server.NewSearchServer(cfg, lockDb, rateLimiter)
	protogen.RegisterSearchServiceServer(s, searchServer)

	// Start a goroutine for daily count reset at UTC-7 midnight
	go server.StartDailyResetTask(cfg, lockDb, searchServer)

	defer system.SafeClose()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
