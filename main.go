package main

import (
	"fmt"
	"html/template"
	"net/http"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", QuoteHandler)

	fmt.Println(http.ListenAndServe(":8080", nil))
}

// QuoteHandler queries the producer and displays the quote
func QuoteHandler(w http.ResponseWriter, r *http.Request) {

	template := template.Must(template.ParseFiles("templates/index.html"))

	var ep Endpoint
	if cfenv.IsRunningOnCF() {
		appEnv, _ := cfenv.Current()
		service, _ := appEnv.Services.WithName("producer-endpoint")
		url := fmt.Sprintf("%v", service.Credentials["url"])
		port := service.Credentials["port"].(int)
		ep = Endpoint{
			URL:    url,
			Port:   port,
			Client: &http.Client{},
		}
	} else {
		ep = Endpoint{
			URL:    "http://localhost:8080",
			Client: &http.Client{},
		}
	}
	quote, err := ep.FetchQuote()
	if err != nil {
		fmt.Printf("error")
	}
	fmt.Printf("Quote: %+v", quote)
	if err := template.ExecuteTemplate(w, "index.html", quote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
