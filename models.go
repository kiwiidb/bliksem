package main

import (
	"net/http"
)

//BliksemService encapsulates the global state of the service
type BliksemService struct {
	client   *http.Client
	url      string
	macaroon string
	invoiceChan chan Invoice
}

//Invoice represents a lightning invoice
type Invoice struct {
	Memo   string `json:"memo"`
	Amount int64  `json:"value"`
	PayReq string `json:"payment_request"`
}

//StringInvoice represents a lightning invoice where all fields are strings
type StringInvoice struct {
	Memo   string `json:"memo"`
	Amount string `json:"value"`
	PayReq string `json:"payment_request"`
}

//LNDStreamInvoice represents an invoice received from LND stream of settled invoices
type LNDStreamInvoice struct {
	Result StringInvoice `json:"result"`
}

//ReqInvoice represents a lightning invoice received from vue frontend
type ReqInvoice struct {
	Body string `json:"body"`
}
