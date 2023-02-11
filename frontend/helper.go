package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kanowfy/ecom/frontend/pb"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (app *application) getProducts(ctx context.Context, query string) ([]*pb.Product, error) {
	resp, err := pb.NewCatalogClient(app.catalogSvcConn).SearchProducts(ctx, &pb.SearchProductsRequest{
		Query: query,
	})

	if err != nil {
		return nil, err
	}

	return resp.Results, nil
}

func (app *application) getProduct(ctx context.Context, id string) (*pb.Product, error) {
	resp, err := pb.NewCatalogClient(app.catalogSvcConn).GetProductById(ctx, &pb.GetProductByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return resp.Product, nil
}

func (app *application) addToCart(ctx context.Context, userid string, productid string, quantity uint32) error {
	_, err := pb.NewCartClient(app.cartSvcConn).AddItem(ctx, &pb.AddItemRequest{
		UserId: userid,
		Item: &pb.CartItem{
			ProductId: productid,
			Quantity:  quantity,
		},
	})

	return err
}

func (app *application) removeFromCart(ctx context.Context, userid string, productid string) error {
	_, err := pb.NewCartClient(app.cartSvcConn).RemoveItem(ctx, &pb.RemoveItemRequest{
		UserId:    userid,
		ProductId: productid,
	})

	return err
}

func (app *application) getCart(ctx context.Context, userid string) ([]*pb.CartItem, error) {
	resp, err := pb.NewCartClient(app.cartSvcConn).GetCart(ctx, &pb.GetCartRequest{
		UserId: userid,
	})

	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func (app *application) getShipCost(ctx context.Context, items []*pb.CartItem) (int64, error) {
	resp, err := pb.NewShippingClient(app.shippingSvcConn).GetShippingCost(ctx, &pb.GetShippingCostRequest{
		Address: nil,
		Items:   items,
	})

	if err != nil {
		return 0, err
	}

	return resp.GetCost(), nil
}

func (app *application) cartDetailsView(w http.ResponseWriter, r *http.Request, page string) {
	cartItems, err := app.getCart(context.Background(), sessionID(r))
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				app.renderTemplates(w, "empty_cart.page.html", nil)
			default:
				log.Println(e.Code(), e.Message())
			}
		} else {
			log.Println(err.Error())
		}
		return
	}

	type itemView struct {
		Item      *pb.Product
		Quantity  uint32
		ItemTotal int64
	}

	items := make([]itemView, len(cartItems))
	var totalItemPrices int64 = 0
	var cart_size uint32 = 0
	for i, item := range cartItems {
		p, err := app.getProduct(context.Background(), item.GetProductId())
		if err != nil {
			log.Println(err.Error())
			return
		}

		items[i] = itemView{
			Item:      p,
			Quantity:  item.Quantity,
			ItemTotal: p.PriceVnd * int64(item.Quantity),
		}

		totalItemPrices += p.PriceVnd * int64(item.Quantity)
		cart_size += item.Quantity
	}

	shipCost, err := app.getShipCost(context.Background(), cartItems)
	if err != nil {
		log.Println(err.Error())
		return
	}

	totalPrice := totalItemPrices + shipCost

	app.renderTemplates(w, page+".page.html", map[string]interface{}{
		"cart_size":                    cart_size,
		"items":                        items,
		"total_price_without_shipping": totalItemPrices,
		"shipping_cost":                shipCost,
		"total_price":                  totalPrice,
	})
}

func (app *application) emptyCart(ctx context.Context, userid string) error {
	cartItems, err := app.getCart(context.Background(), userid)
	if err != nil {
		return fmt.Errorf("could not get cart: %v", err)
	}
	for _, item := range cartItems {
		err := app.removeFromCart(ctx, userid, item.GetProductId())
		if err != nil {
			return err
		}
	}

	return nil
}

func parseDate(date string, part_idx int) uint32 {
	parts := strings.Split(date, "/")
	part, _ := strconv.ParseUint(parts[part_idx], 10, 32)
	return uint32(part)
}

func getMonthPart(date string) uint32 {
	return parseDate(date, 0)
}

func getYearPart(date string) uint32 {
	return parseDate(date, 1)
}

func formattedPrice(price int64) string {
	return message.NewPrinter(language.English).Sprint(price)
}
