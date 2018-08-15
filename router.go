package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const (
	addInvoiceRoute     = "/addinvoice"
	settledInvoiceRoute = "/settledinvoice"
)

func setupRouter(service *BliksemService, conf Config) *negroni.Negroni {
	router := mux.NewRouter()
	router.HandleFunc(addInvoiceRoute,
		func(w http.ResponseWriter, r *http.Request) { handleAddInvoice(w, r, service) })

	router.HandleFunc(settledInvoiceRoute,
		func(w http.ResponseWriter, r *http.Request) { handleSettledInvoice(w, r, service) })

	n := negroni.Classic()

	/*
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
		})
	*/
	c := cors.Default()
	n.Use(c)
	n.UseHandler(router)
	return n
}
