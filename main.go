package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type (
	ResultAPI struct {
		Data DataAPI `json:"data"`
	}
	DataAPI struct {
		Products []Product `json:"products"`
	}

	Product struct {
		ID           string       `json:"id"`
		Name         string       `json:"name"`
		Images       []string     `json:"images"`
		Location     string       `json:"location"`
		Price        Price        `json:"price"`
		RootCategory RootCategory `json:"rootCategory"`
	}

	Price struct {
		PriceDisplay string `json:"priceDisplay"`
	}

	RootCategory struct {
		Name string `json:"name"`
	}
)

const (
	// you need to have your own api, because i use one of Indonesian Market Place just for getting product list by keyword
	yourURL = "https://www.your_url_here.com/test/search/products?page=1&start=0&searchTerm=%s"
)

func main() {
	http.HandleFunc("/get_product", loopHandler)
	log.Println("Listen to 127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

func loopHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	keyword := keys.Get("keyword")

	fmap := template.FuncMap{
		"formatAsTakeOneString": formatAsTakeOneString,
	}

	products, err := getDataProduct(keyword)
	if err != nil {
		log.Println("error when accesing loop html", err)
	}

	t := template.Must(template.New("loop.html").Funcs(fmap).ParseFiles("loop.html"))
	t.Execute(w, &products)
}

func getDataProduct(keyword string) (prod []Product, err error) {
	url := fmt.Sprintf(yourURL, keyword)
	//Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return prod, err
	}

	// header to allow request -> in case you need it
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36`)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return prod, err
	}

	defer resp.Body.Close()

	var record ResultAPI
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		return prod, err
	}

	return record.Data.Products, err
}

func formatAsTakeOneString(valueInOneString []string) string {
	return valueInOneString[0]
}
