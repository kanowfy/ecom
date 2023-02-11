package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/kanowfy/ecom/frontend/pb"
)

func (app *application) renderTemplates(w http.ResponseWriter, filename string, data interface{}) {
	ts, ok := app.templateCache[filename]
	if !ok {
		log.Printf("The template %s does not exit", filename)
		return
	}
	err := ts.Execute(w, data)
	if err != nil {
		log.Print(err)
	}

}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("Session ID: %v", sessionID(r))
	app.renderTemplates(w, "home.page.html", nil)
}

func (app *application) productsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := app.getProducts(context.Background(), "")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.renderTemplates(w, "products.page.html", res)
}

func (app *application) productHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	product, err := app.getProduct(context.Background(), ps.ByName("id"))
	if err != nil {
		log.Println(err.Error())
		return
	}

	app.renderTemplates(w, "product.page.html", product)
}

func (app *application) viewCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.cartDetailsView(w, r, "cart")
}

func (app *application) addToCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	id := r.FormValue("productid")
	quantity, _ := strconv.ParseUint(r.FormValue("quantity"), 10, 32)

	p, err := app.getProduct(context.Background(), id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = app.addToCart(context.Background(), sessionID(r), p.GetId(), uint32(quantity))
	if err != nil {
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
}

func (app *application) removeFromCartHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	id := r.FormValue("productid")

	p, err := app.getProduct(context.Background(), id)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = app.removeFromCart(context.Background(), sessionID(r), p.GetId())
	if err != nil {
		log.Println(err.Error())
		return
	}

	http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
}

func (app *application) checkoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.cartDetailsView(w, r, "checkout")
}

func (app *application) placeOrderHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	resp, err := pb.NewOrderClient(app.orderSvcConn).PlaceOrder(context.Background(), &pb.PlaceOrderRequest{
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
	})
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/checkout", http.StatusMovedPermanently)
		return
	}

	app.renderTemplates(w, "result.page.html", map[string]interface{}{
		"order_id":    resp.Order.GetOrderId(),
		"tracking_id": resp.Order.ShippingTrackingId,
	})

	err = app.emptyCart(context.Background(), sessionID(r))
	if err != nil {
		log.Printf("emptying cart failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "about.page.html", nil)
}

func sessionID(r *http.Request) string {
	id := r.Context().Value(ctxSIDKey{})
	if id != nil {
		return id.(string)
	}

	return ""
}
