package main

import (
	"log"
	"net"

	"github.com/kanowfy/ecom/email_service/mailer"
	"github.com/kanowfy/ecom/email_service/pb"

	"google.golang.org/grpc"
)

var (
	port = ":3005"
)

type mailerConfig struct {
	host     string
	port     int
	username string
	password string
	sender   string
}

type server struct {
	mailer mailer.Mailer
	*pb.UnimplementedEmailServer
}

func main() {
	cfg := mailerConfig{
		host:     "smtp.mailtrap.io",
		port:     2525,
		username: "ad6c3dc5d1ad98",
		password: "b45ea57aabed01",
		sender:   "Ecom <no-reply@ecom.kanowfy.com>",
	}
	srv := &server{
		mailer: mailer.New(cfg.host, cfg.port, cfg.username, cfg.password, cfg.sender),
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterEmailServer(s, srv)
	log.Printf("gRPC server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
