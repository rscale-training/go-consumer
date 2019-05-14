package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Endpoint holds the base config for the producer
type Endpoint struct {
	URL    string
	Port   int
	Client *http.Client
}

// Quote holds the text and author
type Quote struct {
	Text   string `json:"quote"`
	Author string `json:"author"`
}

// FetchQuote retrieves a quote from a producer
func (endpoint *Endpoint) FetchQuote() (*Quote, error) {
	reqURL := endpoint.URL
	if endpoint.Port > 0 {
		reqURL = reqURL + strconv.Itoa(endpoint.Port)
	}
	resp, err := endpoint.Client.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var quote Quote
	json.Unmarshal(data, &quote)
	return &quote, nil
}
