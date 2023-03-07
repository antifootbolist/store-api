package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Product struct {
	Id          int
	Name        string
	Description string
	Price       int
}

type Products struct {
	Products []Product
}

func main() {
	fmt.Println("Starting server on port 8080 ...")

	var err error
	// Establish connection to container DB
	db, err = sql.Open("postgres", "host=postgresql user=user-api password=qwe123 dbname=store_api sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// API handlers
	http.HandleFunc("/api/v1/products", GetProducts)
	http.HandleFunc("/api/v1/product/", GetUpdateDeleteProduct)

	// Run server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetUpdateDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// TODO: Add DELETE method
	if r.Method == "GET" {
		GetProduct(w, r)
	} else if r.Method == "POST" {
		UpdateProduct(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

/**
* @api {get} /products Get list of products
 * @apiName GetProducts
 * @apiGroup Products
 * @apiPermission user
 *
 * @apiSuccess {Array} products List of products.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "products": [
 *         {
 *           "Id": 1,
 *           "Name": "Product 1",
 *           "Description": "Product 1 description",
 *           "Price": 10
 *         },
 *         {
 *           "Id": 2,
 *           "Name": "Product 2",
 *           "Description": "Product 2 description",
 *           "Price": 20
 *         },
 *         {
 *           "Id": 3,
 *           "Name": "Product 3",
 *           "Description": "Product 3 description",
 *           "Price": 30
 *         }
 *       ]
 *     }
 *
 * @apiError InvalidRequest The request is invalid.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 400 Bad Request
 *     {
 *       "error": "InvalidRequest"
 *     }
*/

// func GetProducts(w http.ResponseWriter, r *http.Request) {
func GetProducts(w http.ResponseWriter, r *http.Request) {

	w_products := Products{}

	// Query to DB to get list of products
	fmt.Println("# Query from table products")
	rows, err := db.Query("SELECT id,name,description,price FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		w_product := Product{}

		err = rows.Scan(&w_product.Id, &w_product.Name, &w_product.Description, &w_product.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w_products.Products = append(w_products.Products, w_product)
	}

	// Transform to JSON format and send to a client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(w_products)
}

/**
* @api {get} /product/:id Get Product's Info
 * @apiName GetProduct
 * @apiGroup Products
 * @apiPermission user
 *
 * @apiParam {Number} id Product ID.
 *
 * @apiSuccess {Number} id Product ID.
 * @apiSuccess {String} name Name of the Product.
 * @apiSuccess {String} description Description of the Product.
 * @apiSuccess {Number} price Price of the Product.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "product": {
 *         "Id": 123,
 *         "Name": "Product name",
 *         "Description": "Product description",
 *         "Price": 100
 *       }
 *     }
 *
 * @apiError ProductNotFound The product was not found.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 404 Not Found
 *     {
 *       "error": "ProductNotFound"
 *     }
*/

func GetProduct(w http.ResponseWriter, r *http.Request) {

	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to get the product with the given ID
	row := db.QueryRow("SELECT id,name,description,price FROM products WHERE id = $1", id)

	// Initialize a new product
	product := Product{}

	// Fill the product struct with data from the row
	err = row.Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err == sql.ErrNoRows {
		http.Error(w, "ProductNotFound", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Transform to JSON format and send to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
}

/**
* @api {post} /product/:id Update Product's Info
 * @apiName UpdateProduct
 * @apiGroup Products
 * @apiPermission admin
 *
 * @apiParam {Number} id Product ID.
 *
 * @apiParam {String} [name] Name of the Product.
 * @apiParam {String} [description] Description of the Product.
 * @apiParam {Number} [price] Price of the Product.
 *
 * @apiSuccess {Object} product Updated product information.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "product": {
 *         "Id": 123,
 *         "Name": "New name",
 *         "Description": "New description",
 *         "Price": 100
 *       }
 *     }
 *
 * @apiError ProductNotFound The product was not found.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 404 Not Found
 *     {
 *       "error": "ProductNotFound"
 *     }
 *
*/

func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the product information from the request body
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build the SQL query to update the product
	query := "UPDATE products SET "
	var values []interface{}
	var index = 1

	if updatedProduct.Name != "" {
		query += fmt.Sprintf("name = $%v, ", index)
		values = append(values, updatedProduct.Name)
		index++
	}

	if updatedProduct.Description != "" {
		query += fmt.Sprintf("description = $%v, ", index)
		values = append(values, updatedProduct.Description)
		index++
	}

	if updatedProduct.Price != 0 {
		query += fmt.Sprintf("price = $%d, ", index)
		values = append(values, updatedProduct.Price)
		index++
	}

	// Form the WHERE clause based on the number of values
	whereIndex := len(values) + 1
	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", whereIndex)

	// Update the product in the database
	fmt.Println("# Update product in table products")
	fmt.Println("### Query", query)
	fmt.Println("### Values", values)
	result, err := db.Exec(query, append(values, id)...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no rows were affected, the product was not found
	if rowsAffected == 0 {
		http.Error(w, "ProductNotFound", http.StatusNotFound)
		return
	}

	// Query the database to get the updated product information
	fmt.Println("# Query from table products")
	rows, err := db.Query("SELECT id,name,description,price FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the updated product information into a Products struct
	w_products := Products{}
	for rows.Next() {
		w_product := Product{}
		err = rows.Scan(&w_product.Id, &w_product.Name, &w_product.Description, &w_product.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w_products.Products = append(w_products.Products, w_product)
	}

	// Return the updated product information to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"product": w_products.Products[0]})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
