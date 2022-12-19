package main

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/kanowfy/ecom/order_service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *server) PlaceOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	log.Printf("received a PlaceOrder request with data: %+v", req)
	// prepare order id
	id, err := uuid.NewV4()
	if err != nil {
		log.Printf("can not generate uuid: %v", err)
		return nil, status.Errorf(codes.Internal, "unable to generate uuid: %v", err)
	}

	// get user cart
	cart, err := pb.NewCartClient(s.conns.cartSvcConn).GetCart(ctx, &pb.GetCartRequest{UserId: req.GetUserId()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: unable to get user cart during checkout: %+v", err)
	}

	cartItems := cart.GetItems()

	// infer order items from user cart
	orderItems := make([]*pb.OrderItem, len(cartItems))
	catalogClient := pb.NewCatalogClient(s.conns.catalogSvcConn)
	for i, item := range cartItems {
		product, err := catalogClient.GetProductById(ctx, &pb.GetProductByIdRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to prepare order: unable to get product #%q", item.GetProductId())
		}
		orderItems[i] = &pb.OrderItem{
			Item: item,
			Cost: product.GetProduct().GetPriceVnd(),
		}
	}

	// get shipping cost
	shippingVND, err := pb.NewShippingClient(s.conns.shippingSvcConn).GetShippingCost(ctx, &pb.GetShippingCostRequest{
		Address: req.GetAddress(),
		Items:   cartItems,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get shipping cost: %+v", err)
	}

	// calculate total cost
	var totalCost int64 = 0
	for _, item := range orderItems {
		totalCost += item.GetCost()
	}

	totalCost += shippingVND.GetCost()

	// charge client
	paymentResp, err := pb.NewPaymentClient(s.conns.paymentSvcConn).Charge(ctx, &pb.ChargeRequest{
		Amount:     totalCost,
		CreditCard: req.GetCard(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}

	log.Printf("payment went through (transaction id: %s)", paymentResp.PaymentId)

	// ship order
	shipResp, err := pb.NewShippingClient(s.conns.shippingSvcConn).ShipOrder(ctx, &pb.ShipOrderRequest{
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
	_, err = pb.NewEmailClient(s.conns.emailSvcConn).SendConfirmation(ctx, &pb.SendConfirmationRequest{
		Email: req.GetEmail(),
		Order: orderResults,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send confirmation email: %+v", err)
	}

	log.Printf("confirmation email sent to %q", req.GetEmail())

	return &pb.PlaceOrderResponse{Order: orderResults}, nil

}
