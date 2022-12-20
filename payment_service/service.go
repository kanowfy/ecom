package main

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/payment_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	log.Printf("received a Charge request with data: %v", req)
	err := VerifyCard(req.GetCreditCard().CardNumber, req.GetCreditCard().CardCvv, req.GetCreditCard().CardExpirationYear, req.GetCreditCard().CardExpirationMonth)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "card verification failed: %v", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("can not generate uuid: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to generate uuid: %v", err)
	}

	log.Println("Charge successful")
	resp := &pb.ChargeResponse{
		PaymentId: id.String(),
	}
	return resp, status.New(codes.OK, "").Err()
}
