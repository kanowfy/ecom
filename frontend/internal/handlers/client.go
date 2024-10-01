package handlers

import (
	"context"
	"fmt"

	"github.com/kanowfy/ecom/frontend/pb"
)

func (h *Handlers) getProducts(ctx context.Context, query string) ([]*pb.Product, error) {
	resp, err := pb.NewCatalogClient(h.conns.CatalogSvc).SearchProducts(ctx, &pb.SearchProductsRequest{
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (h *Handlers) getProduct(ctx context.Context, id string) (*pb.Product, error) {
	resp, err := pb.NewCatalogClient(h.conns.CatalogSvc).GetProductById(ctx, &pb.GetProductByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}

func (h *Handlers) addToCart(ctx context.Context, userid string, productid string, quantity uint32) error {
	_, err := pb.NewCartClient(h.conns.CartSvc).AddItem(ctx, &pb.AddItemRequest{
		UserId: userid,
		Item: &pb.CartItem{
			ProductId: productid,
			Quantity:  quantity,
		},
	})

	return err
}

func (h *Handlers) removeFromCart(ctx context.Context, userid string, productid string) error {
	_, err := pb.NewCartClient(h.conns.CartSvc).RemoveItem(ctx, &pb.RemoveItemRequest{
		UserId:    userid,
		ProductId: productid,
	})

	return err
}

func (h *Handlers) getCart(ctx context.Context, userid string) ([]*pb.CartItem, error) {
	resp, err := pb.NewCartClient(h.conns.CartSvc).GetCart(ctx, &pb.GetCartRequest{
		UserId: userid,
	})

	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (h *Handlers) getShipCost(ctx context.Context, items []*pb.CartItem) (int64, error) {
	resp, err := pb.NewShippingClient(h.conns.ShippingSvc).GetShippingCost(ctx, &pb.GetShippingCostRequest{
		Address: nil,
		Items:   items,
	})

	if err != nil {
		return 0, err
	}

	return resp.GetCost(), nil
}

func (h *Handlers) emptyCart(ctx context.Context, userid string) error {
	cartItems, err := h.getCart(ctx, userid)
	if err != nil {
		return fmt.Errorf("could not get cart: %v", err)
	}
	for _, item := range cartItems {
		err := h.removeFromCart(ctx, userid, item.GetProductId())
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handlers) placeOrder(ctx context.Context, req *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	return pb.NewOrderClient(h.conns.OrderSvc).PlaceOrder(ctx, req)
}
