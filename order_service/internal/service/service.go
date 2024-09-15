package service

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/order_service/internal/grpcconn"
	"github.com/kanowfy/ecom/order_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	logger *slog.Logger
	conns  *grpcconn.Connection
	*pb.UnimplementedOrderServer
}

func New(logger *slog.Logger, conns *grpcconn.Connection) *service {
	return &service{
		logger: logger,
		conns:  conns,
	}
}

func (s *service) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	s.logger.Info("received a PlaceOrder request")
	// prepare order id
	id, err := uuid.NewV4()
	if err != nil {
		s.logger.Error("failed to generate uuid", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "unable to generate uuid: %v", err)
	}

	logger := s.logger.With(slog.String("user_id", req.GetUserId()), slog.String("order_id", id.String()))

	// get user cart
	cart, err := pb.NewCartClient(s.conns.CartSvc).GetCart(ctx, &pb.GetCartRequest{UserId: req.GetUserId()})
	if err != nil {
		logger.Error("failed to get user cart", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get cart: unable to get user cart during checkout: %+v", err)
	}

	cartItems := cart.GetItems()

	// infer order items from user cart
	orderItems := make([]*pb.OrderItem, len(cartItems))
	catalogClient := pb.NewCatalogClient(s.conns.CatalogSvc)
	for i, item := range cartItems {
		product, err := catalogClient.GetProductById(ctx, &pb.GetProductByIdRequest{Id: item.GetProductId()})
		if err != nil {
			logger.Error("failed to fetch product", "error", err.Error(),
				"product_id", item.GetProductId())
			return nil, status.Errorf(codes.Internal, "failed to prepare order: unable to get product #%q", item.GetProductId())
		}
		orderItems[i] = &pb.OrderItem{
			Item: item,
			Cost: product.GetProduct().GetPriceVnd(),
		}
	}

	// get shipping cost
	shippingVND, err := pb.NewShippingClient(s.conns.ShippingSvc).GetShippingCost(ctx, &pb.GetShippingCostRequest{
		Address: req.GetAddress(),
		Items:   cartItems,
	})
	if err != nil {
		logger.Error("failed to get shipping cost", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to get shipping cost: %+v", err)
	}

	// calculate total cost
	var totalCost int64 = 0
	for _, item := range orderItems {
		totalCost += item.GetCost()
	}

	totalCost += shippingVND.GetCost()

	// charge client
	paymentResp, err := pb.NewPaymentClient(s.conns.PaymentSvc).Charge(ctx, &pb.ChargeRequest{
		Amount:     totalCost,
		CreditCard: req.GetCard(),
	})
	if err != nil {
		logger.Error("failed to charge", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}

	logger.Info("payment succeeded", "payment_id", paymentResp.PaymentId)

	// ship order
	shipResp, err := pb.NewShippingClient(s.conns.ShippingSvc).ShipOrder(ctx, &pb.ShipOrderRequest{
		Address: req.GetAddress(),
		Items:   cartItems,
	})

	// prepare order results
	orderResults := &pb.OrderResult{
		OrderId:            id.String(),
		ShippingTrackingId: shipResp.GetTrackingId(),
		ShippingCost:       shippingVND.GetCost(),
		ShippingAddress:    req.GetAddress(),
		Items:              orderItems,
	}

	// send confirmation email
	_, err = pb.NewEmailClient(s.conns.EmailSvc).SendConfirmation(ctx, &pb.SendConfirmationRequest{
		Email: req.GetEmail(),
		Order: orderResults,
	})

	if err != nil {
		logger.Error("failed to send confirmation email", "error", err.Error())
		return nil, status.Errorf(codes.Internal, "failed to send confirmation email: %+v", err)
	}

	logger.Info("PlaceOrder successful")

	return &pb.PlaceOrderResponse{Order: orderResults}, nil

}
