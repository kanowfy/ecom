package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/kanowfy/ecom/shipping_service/pb"
)

type service struct {
	logger *slog.Logger
	*pb.UnimplementedShippingServer
}

func New(logger *slog.Logger) *service {
	return &service{
		logger: logger,
	}
}

func (s *service) GetShippingCost(ctx context.Context, req *pb.GetShippingCostRequest) (*pb.GetShippingCostResponse, error) {
	s.logger.Info("received a GetShippingCost request")
	defer s.logger.Info("GetShippingCost successful")
	return &pb.GetShippingCostResponse{
		Cost: 50000,
	}, nil
}

func (s *service) ShipOrder(ctx context.Context, req *pb.ShipOrderRequest) (*pb.ShipOrderResponse, error) {
	s.logger.Info("received a ShipOrder request")
	defer s.logger.Info("ShipOrder created successfully")

	trackingId := createTrackingId()

	return &pb.ShipOrderResponse{
		TrackingId: trackingId,
	}, nil
}

func createTrackingId() string {
	return fmt.Sprintf("%c%c-%s-%s",
		generateRandomCharacter(),
		generateRandomCharacter(),
		generateRandomNumber(5),
		generateRandomNumber(9),
	)
}

func generateRandomCharacter() uint32 {
	return 65 + uint32(rand.Intn(25))
}

func generateRandomNumber(length int) string {
	s := ""
	for i := 0; i < length; i++ {
		s = fmt.Sprintf("%s%d", s, rand.Intn(10))
	}

	return s
}
