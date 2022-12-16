package main

import (
	"github.com/kanowfy/ecom/shipping_service/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = ":3004"
)

type server struct {
	*pb.UnimplementedShippingServer
}

func main() {
	srv := &server{}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterShippingServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
