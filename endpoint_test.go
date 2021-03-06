package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchQuote(t *testing.T) {

	expectedBytes := []byte(`{"quote":"This is a quote","author":"Quotey McQuoteface"}`)

	expectedQuote := Quote{
		Text:   "This is a quote",
		Author: "Quotey McQuoteface",
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(expectedBytes)
	}))
	defer server.Close()

	endpoint := Endpoint{URL: server.URL, Client: server.Client()}
	quote, err := endpoint.FetchQuote()
	if err != nil {
		t.Errorf("Error fetching quote: %e", err)
	}

	if quote.Text != expectedQuote.Text || quote.Author != expectedQuote.Author {
		t.Errorf("Quote does not match. Expected: %+v. Received: %+v.", expectedQuote, quote)
	}
}

func TestGetPortFromString(t *testing.T) {
	var i interface{} = "8080"
	port, _ := GetPort(i)
	if port != 8080 {
		t.Errorf("Parsing port from string failed")
	}
}

func TestGetPortFromFloat64(t *testing.T) {
	var i interface{} = float64(8080)
	port, _ := GetPort(i)
	if port != 8080 {
		t.Errorf("Parsing port from string failed")
	}
}
