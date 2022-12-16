package main

import (
	"flag"
	"github.com/kanowfy/ecom/email_service/mailer"
	"github.com/kanowfy/ecom/email_service/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = ":3005"
)

type mailerConfig struct {
	host string
	port int
	username string
	password string
	sender string
}

type server struct {
	mailer mailer.Mailer
	*pb.UnimplementedEmailServer
}

func main() {
	var cfg mailerConfig
	flag.StringVar(&cfg.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.username, "smtp-username", "smtp username (to be filled)", "SMTP username")
	flag.StringVar(&cfg.password, "smtp-password", "smtp password (to be filled)", "SMTP password")
	flag.StringVar(&cfg.sender, "smtp-sender", "Ecom <no-reply@ecom.kanowfy.com>", "SMTP sender")
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
