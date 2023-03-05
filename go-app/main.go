package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

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
	fmt.Println("Starting server on port 8080 ...")

	// Handler API
	http.HandleFunc("/api/v1/products", getProducts)

	// Run server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w_products := Products{}

	// Establish connection to container DB
	db, err := sql.Open("postgres", "host=130.193.36.79 user=user-api password=qwe123 dbname=store_api sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Query to DB to get list of products
	fmt.Println("# Query from table products")
	rows, err := db.Query("SELECT id,name,description,price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("# Value of rows", rows)
	for rows.Next() {
		w_product := Product{}

		err = rows.Scan(&w_product.id, &w_product.name, &w_product.description, &w_product.price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w_products.Products = append(w_products.Products, w_product)
	}

	fmt.Println("# Value of Products", w_products)
	// Transform to JSON format and send to client
	json.NewEncoder(w).Encode(w_products)

}
