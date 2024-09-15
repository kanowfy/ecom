package service

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/payment_service/internal/validator"
	"github.com/kanowfy/ecom/payment_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	logger *slog.Logger
	*pb.UnimplementedPaymentServer
}

func New(logger *slog.Logger) *service {
	return &service{
		logger: logger,
	}
}

func (s *service) Charge(ctx context.Context, req *pb.ChargeRequest) (*pb.ChargeResponse, error) {
	s.logger.Info("received a Charge request")

	err := validator.VerifyCard(req.GetCreditCard().CardNumber, req.GetCreditCard().CardCvv, req.GetCreditCard().CardExpirationYear, req.GetCreditCard().CardExpirationMonth)
	if err != nil {
		s.logger.Error("card validation failed", "error", err.Error())
		return nil, status.Errorf(codes.InvalidArgument, "card verification failed: %v", err)
	}

	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed to generate uuid", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "unable to generate uuid: %v", err)
	}

	s.logger.Info("Charge successful")
	resp := &pb.ChargeResponse{
		PaymentId: id.String(),
	}
	return resp, nil
}
