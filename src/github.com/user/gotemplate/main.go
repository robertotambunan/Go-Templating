package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type ReturnAPI struct {
	Product []Product `json:"products"`
}
type Product struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	City           string   `json:"city"`
	Price          int      `json:"price"`
	Category       string   `json:"category"`
	SellerUsername string   `json:"seller_username"`
	SellerName     string   `json:"seller_name"`
	Province       string   `json:"province"`
	Url            string   `json:"url"`
	Weight         int      `json:"weight"`
	Stock          int      `json:"stock"`
	Images         []string `json:"images"`
	Small_Images   []string `json:"small_images"`
	Urls           string   `json:"url"`
}

func main() {
	http.HandleFunc("/", loopHandler)
	log.Println("Listen to 127.0.0.1:8080")
	http.ListenAndServe(":8080", nil)
}

func loopHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Path[1:]
	fmap := template.FuncMap{
		"formatAsInt":           formatAsInt,
		"formatAsMoney":         formatAsMoney,
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
	url := fmt.Sprintf("https://api.bukalapak.com/v2/products.json?keywords=%s&page=1&per_page=5", keyword)
	//Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return prod, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return prod, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var record ReturnAPI
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
		return prod, err
	}

	return record.Product, err
}

func formatAsInt(integerForm int) (numbers string) {
	numbers = strconv.Itoa(integerForm)
	return
}

func formatAsMoney(valueInCents int) (string, error) {
	dollars := valueInCents
	return fmt.Sprintf("Rp. %d,00", dollars), nil
}

func formatAsTakeOneString(valueInOneString []string) string {
	return valueInOneString[0]
}
