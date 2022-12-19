package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/kanowfy/ecom/cart_service/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var (
	port = ":3002"
)

type server struct {
	model CartModel
	*pb.UnimplementedCartServer
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RD_URL"),
		Password: "",
		DB:       0,
	})

	srv := &server{
		model: NewModel(rdb),
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterCartServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
