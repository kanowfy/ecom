package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kanowfy/ecom/catalog_service/pb"
	"google.golang.org/grpc"
)

var (
	port = ":3001"
)

type server struct {
	model ProductModel
	*pb.UnimplementedCatalogServer
}

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	srv := &server{
		model: NewModel(pool),
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterCatalogServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
