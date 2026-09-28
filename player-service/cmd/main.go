package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	playerpb "study/api"
	"study/internal/handler"
	"study/internal/postgres"
	"study/internal/usecase"
)

func main() {
	_ = godotenv.Load()

	httpAddr := flag.String("http-addr",
		getEnv("HTTP_ADDR", ":8080"), "HTTP listen address")
	grpcAddr := flag.String("grpc-addr",
		getEnv("GRPC_ADDR", ":5600"), "gRPC listen address")
	flag.Parse()

	initCtx, initCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer initCancel()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL is required")
	}

	pool, err := pgxpool.New(initCtx, connStr)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(initCtx); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("db connected")

	r := postgres.NewRepo(pool)

	if err := r.CreateTable(initCtx); err != nil {
		log.Fatalf("create table: %v", err)
	}
	log.Println("table ready")

	uc := usecase.New(r)

	hHTTP := handler.NewHTTP(uc)
	httpSrv := &http.Server{
		Addr:    *httpAddr,
		Handler: hHTTP.Router(),
	}

	hGRPC := handler.NewGRPC(uc)
	grpcSrv := grpc.NewServer()
	playerpb.RegisterPlayerServiceServer(grpcSrv, hGRPC)

	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("grpc listen: %v", err)
	}

	go func() {
		log.Printf("HTTP on %s", *httpAddr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	go func() {
		log.Printf("gRPC on %s", *grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	sigCtx, sigCancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer sigCancel()

	<-sigCtx.Done()
	log.Println("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	grpcSrv.GracefulStop()
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	log.Println("stopped")
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
