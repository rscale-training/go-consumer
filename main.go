package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

// ErrorPage holds the error message
type ErrorPage struct {
	Message string
}

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", QuoteHandler)

	fmt.Println(http.ListenAndServe(":8080", nil))
}

// QuoteHandler queries the producer and displays the quote
func QuoteHandler(w http.ResponseWriter, r *http.Request) {

	template := template.Must(template.ParseFiles("templates/index.html", "templates/error.html"))

	var ep Endpoint
	if cfenv.IsRunningOnCF() {
		appEnv, _ := cfenv.Current()
		service, err := appEnv.Services.WithName("producer-endpoint")

		if err != nil {
			fmt.Println("WARNING: User Provided Service Instance named `producer-endpoint` with `url` and `port` credentials is missing")
			errorPage := ErrorPage{"Error contacting the producer at url: undefined on port: undefined"}
			if err := template.ExecuteTemplate(w, "error.html", errorPage); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		url := fmt.Sprintf("%v", service.Credentials["url"])
		url = "http://" + url
		port, _ := GetPort(service.Credentials["port"])

		ep = Endpoint{
			URL:    url,
			Port:   port,
			Client: &http.Client{},
		}
	} else {
		ep = Endpoint{
			URL:    "http://localhost",
			Port:   8080,
			Client: &http.Client{},
		}
	}
	quote, err := ep.FetchQuote()
	if err != nil {
		errorMessage := "ERROR contacting the producer at url: " + ep.URL + " on port: " + strconv.Itoa(ep.Port)
		fmt.Println(errorMessage)
		errorPage := ErrorPage{errorMessage}
		if err := template.ExecuteTemplate(w, "error.html", errorPage); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := template.ExecuteTemplate(w, "index.html", quote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
