package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kanowfy/ecom/frontend/internal/config"
	"github.com/kanowfy/ecom/frontend/internal/handlers"
	"github.com/kanowfy/ecom/frontend/internal/middleware"
)

func New(h *handlers.Handlers, cfg config.Config) http.Handler {
	r := httprouter.New()
	r.GET("/", h.HomeHandler)
	r.GET("/products", h.ProductsHandler)
	r.GET("/products/:id", h.ProductHandler)
	r.GET("/cart", h.ViewCartHandler)
	r.POST("/cart", h.AddToCartHandler)
	r.POST("/cart/remove", h.RemoveFromCartHandler)
	r.GET("/checkout", h.CheckoutHandler)
	r.POST("/checkout", h.PlaceOrderHandler)
	r.GET("/about", h.AboutHandler)
	r.ServeFiles("/public/*filepath", http.Dir("./ui/static"))

	var handler http.Handler = r
	handler = middleware.CheckSessionID(handler, cfg.Cookie.SID, cfg.Cookie.MaxAge)

	return handler
}
