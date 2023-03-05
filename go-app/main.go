package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	id          int
	name        string
	description string
	price       int
}

type Products struct {
	Products []Product
}

func main() {
	var err error

	// Establish connection to container DB
	db, err := sql.Open("postgres", "host=postgresql user=user-api password=qwe123 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Handler API
	http.HandleFunc("/api/v1/products", getProducts)

	// Run server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w_products := Products{}

	// Query to DB to get list of products
	rows, err := db.Query("SELECT id, name, description, price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		w_product := Product{}

		err = rows.Scan(&w_product.id, &w_product.name, &w_product.description, &w_product.price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w_products.Products = append(w_products.Products, w_product)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform to JSON format and send to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(w_products)

}
