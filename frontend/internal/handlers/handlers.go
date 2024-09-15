package handlers

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kanowfy/ecom/frontend/internal/grpcconn"
	"github.com/kanowfy/ecom/frontend/internal/middleware"
	"github.com/kanowfy/ecom/frontend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handlers struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
	conns         *grpcconn.Connection
}

func New(logger *slog.Logger, templateCache map[string]*template.Template, conns *grpcconn.Connection) *Handlers {
	return &Handlers{
		logger,
		templateCache,
		conns,
	}
}

func (h *Handlers) renderTemplates(w http.ResponseWriter, filename string, data interface{}) {
	ts, ok := h.templateCache[filename]
	if !ok {
		log.Printf("The template %s does not exit", filename)
		return
	}
	err := ts.Execute(w, data)
	if err != nil {
		log.Print(err)
	}

}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("Session ID: %v", sessionID(r))
	h.renderTemplates(w, "home.page.html", nil)
}

func (h *Handlers) ProductsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.getProducts(r.Context(), "")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.renderTemplates(w, "products.page.html", res)
}

func (h *Handlers) ProductHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	product, err := h.getProduct(r.Context(), ps.ByName("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	h.renderTemplates(w, "product.page.html", product)
}

func (h *Handlers) cartDetailsView(w http.ResponseWriter, r *http.Request, page string) {
	cartItems, err := h.getCart(r.Context(), sessionID(r))
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				h.renderTemplates(w, "empty_cart.page.html", nil)
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
		p, err := h.getProduct(r.Context(), item.GetProductId())
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

	shipCost, err := h.getShipCost(r.Context(), cartItems)
	if err != nil {
		log.Println(err.Error())
		return
	}

	totalPrice := totalItemPrices + shipCost

	h.renderTemplates(w, page+".page.html", map[string]interface{}{
		"cart_size":                    cart_size,
		"items":                        items,
		"total_price_without_shipping": totalItemPrices,
		"shipping_cost":                shipCost,
		"total_price":                  totalPrice,
	})
}

func (h *Handlers) ViewCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.cartDetailsView(w, r, "cart")
}

func (h *Handlers) AddToCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	id := r.FormValue("productid")
	quantity, _ := strconv.ParseUint(r.FormValue("quantity"), 10, 32)

	p, err := h.getProduct(r.Context(), id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = h.addToCart(r.Context(), sessionID(r), p.GetId(), uint32(quantity))
	if err != nil {
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
}

func (h *Handlers) RemoveFromCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	id := r.FormValue("productid")

	p, err := h.getProduct(r.Context(), id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = h.removeFromCart(r.Context(), sessionID(r), p.GetId())
	if err != nil {
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
}

func (h *Handlers) CheckoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h.cartDetailsView(w, r, "checkout")
}

func (h *Handlers) PlaceOrderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var (
		email       = r.FormValue("email")
		address     = r.FormValue("address")
		city        = r.FormValue("city")
		country     = r.FormValue("country")
		zip_code, _ = strconv.ParseInt(r.FormValue("zip_code"), 10, 32)
		card_number = r.FormValue("card_number")
		card_exp    = r.FormValue("card_exp")
		card_cvv, _ = strconv.ParseUint(r.FormValue("card_cvv"), 10, 32)
	)

	card_exp_month := getMonthPart(card_exp)
	card_exp_year := getYearPart(card_exp)

	req := &pb.PlaceOrderRequest{
		UserId: sessionID(r),
		Address: &pb.Address{
			StreetAddress: address,
			City:          city,
			Country:       country,
			ZipCode:       int32(zip_code),
		},
		Email: email,
		Card: &pb.CardInfo{
			CardNumber:          card_number,
			CardExpirationMonth: card_exp_month,
			CardExpirationYear:  card_exp_year,
			CardCvv:             uint32(card_cvv),
		},
	}

	resp, err := h.placeOrder(r.Context(), req)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/checkout", http.StatusMovedPermanently)
		return
	}

	h.renderTemplates(w, "result.page.html", map[string]interface{}{
		"order_id":    resp.Order.GetOrderId(),
		"tracking_id": resp.Order.ShippingTrackingId,
	})

	err = h.emptyCart(r.Context(), sessionID(r))
	if err != nil {
		log.Printf("emptying cart failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) AboutHandler(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	h.renderTemplates(w, "about.page.html", nil)
}

func sessionID(r *http.Request) string {
	id := r.Context().Value(middleware.CtxSIDKey{})
	if id != nil {
		return id.(string)
	}

	return ""
}
