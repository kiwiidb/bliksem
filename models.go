package main

import (
	"net/http"
)

//BliksemService encapsulates the global state of the service
type BliksemService struct {
	client   *http.Client
	url      string
	macaroon string
}

//Invoice represents a lightning invoice
type Invoice struct {
	Memo   string `json:"memo"`
	Amount int    `json:"value"`
	PayReq string `json:"payment_request"`
}
