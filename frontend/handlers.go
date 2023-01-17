package main

import (
	"context"
	"log"
	"net/http"

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

func (app *application) home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "home.page.html", nil)
}

func (app *application) products(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp, err := pb.NewCatalogClient(app.catalogSvcConn).SearchProducts(context.Background(), &pb.SearchProductsRequest{
		Query: "",
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.renderTemplates(w, "products.page.html", resp.Results)
}

func (app *application) product(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp, err := pb.NewCatalogClient(app.catalogSvcConn).GetProductById(context.Background(), &pb.GetProductByIdRequest{
		Id: ps.ByName("id"),
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	app.renderTemplates(w, "product.page.html", resp.Product)
}

func (app *application) cart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "cart.page.html", nil)
}

func (app *application) checkout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "checkout.page.html", nil)
}

func (app *application) about(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "about.page.html", nil)
}

func (app *application) result(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.renderTemplates(w, "result.page.html", nil)
}
