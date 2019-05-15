package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		reqURL = reqURL + ":" + strconv.Itoa(endpoint.Port)
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
	fmt.Printf("DATA: %v \n", string(data))
	var quote Quote
	json.Unmarshal(data, &quote)
	if quote.Author == "" || quote.Text == "" {
		return nil, errors.New("Error fetching quote from producer")
	}
	return &quote, nil
}

// GetPort converts a string or float64 (the two types provided by CF) to an int
func GetPort(port interface{}) (int, error) {
	s, ok := port.(string)
	if ok {
		return strconv.Atoi(s)
	}
	f, ok := port.(float64)
	if ok {
		return int(f), nil
	}
	return 80, nil
}
