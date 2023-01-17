package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

var (
	port = ":4000"
)

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
	r.GET("/", app.home)
	r.GET("/products", app.products)
	r.GET("/products/:id", app.product)
	r.GET("/cart", app.cart)
	r.GET("/checkout", app.checkout)
	r.GET("/result", app.result)
	r.GET("/about", app.result)
	r.ServeFiles("/public/*filepath", http.Dir("./ui/static"))

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
