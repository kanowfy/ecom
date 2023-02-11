package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

var (
	port         = ":4000"
	cookieSID    = "ecom_sid"
	cookieMaxAge = 60 * 60 * 24
)

type ctxSIDKey struct{}

type application struct {
	templateCache map[string]*template.Template

	catalogSvcAddr string
	catalogSvcConn *grpc.ClientConn

	cartSvcAddr string
	cartSvcConn *grpc.ClientConn

	shippingSvcAddr string
	shippingSvcConn *grpc.ClientConn

	emailSvcAddr string
	emailSvcConn *grpc.ClientConn

	paymentSvcAddr string
	paymentSvcConn *grpc.ClientConn

	orderSvcAddr string
	orderSvcConn *grpc.ClientConn
}

func main() {
	var app application
	mapConnections(&app)

	templateCache, err := newTemplateCache("./ui/templates/")
	if err != nil {
		log.Fatal(err)
	}

	app.templateCache = templateCache

	r := httprouter.New()
	r.GET("/", app.homeHandler)
	r.GET("/products", app.productsHandler)
	r.GET("/products/:id", app.productHandler)
	r.GET("/cart", app.viewCartHandler)
	r.POST("/cart", app.addToCartHandler)
	r.POST("/cart/remove", app.removeFromCartHandler)
	r.GET("/checkout", app.checkoutHandler)
	r.POST("/checkout", app.placeOrderHandler)
	r.GET("/about", app.aboutHandler)
	r.ServeFiles("/public/*filepath", http.Dir("./ui/static"))

	var handler http.Handler = r
	handler = checkSessionID(handler)

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
