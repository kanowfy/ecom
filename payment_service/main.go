package main

import (
	"log"
	"net"

	"github.com/kanowfy/ecom/payment_service/pb"
	"google.golang.org/grpc"
)

var (
	port = ":3003"
)

type server struct {
	*pb.UnimplementedPaymentServer
}

func main() {
	srv := &server{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterPaymentServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
