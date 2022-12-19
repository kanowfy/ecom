package main

import (
	"github.com/kanowfy/ecom/order_service/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = ":3003"
)

type connection struct {
	catalogSvcAddr string
	catalogSvcConn *grpc.ClientConn

	cartSvcAddr string
	cartSvcConn *grpc.ClientConn

	shippingSvcAddr string
	shippingSvcConn *grpc.ClientConn

	emailSvcAddr string
	emailSvcConn *grpc.ClientConn

	paymentSvcAddr string
	paymentSvcConn *grpc.ClientConn
}

type server struct {
	conns connection
	*pb.UnimplementedOrderServer
}

func main() {
	srv := &server{}
	mapConnections(srv)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
