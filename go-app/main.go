package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Product struct {
	Id          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type Products struct {
	Products []Product `json:"products"`
}

func main() {
	fmt.Println("Starting server on port 8080 ...")

	// Establish connection to container DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	migrate(db)

	// API handlers
	http.HandleFunc("/api/v1/product", GetProducts)
	http.HandleFunc("/api/v1/product/list/", ListProduct)
	http.HandleFunc("/api/v1/product/update/", UpdateProduct)

	// Run server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func migrate(db *gorm.DB) {
	// AutoMigrate will automatically create the table based on the struct
	db.AutoMigrate(&Product{})

	// If TEST_DATA is set to true, insert test data
	testData, err := strconv.ParseBool(os.Getenv("TEST_DATA"))
	if err != nil {
		testData = false
	}
	if testData {
		fmt.Println("Inserting test data...")
		db.Create(&Product{Id: 1, Name: "iPhone", Description: "iPhone 14", Price: 100})
		db.Create(&Product{Id: 2, Name: "iPhone", Description: "iPhone 14 PRO MAX", Price: 200})
		db.Create(&Product{Id: 3, Name: "Samsung", Description: "Samsung Galaxy S23 Ultra", Price: 300})
		fmt.Println("Test data inserted.")
	}
}

/**
* @api {get} /product Get list of all products
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
 *
 * @apiError (Server Error 5xx) InternalServerError There was a server-side error while processing the request.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *       "error": "InternalServerError"
 *     }
*/

func GetProducts(w http.ResponseWriter, r *http.Request) {
	// Check that method is allowed
	if r.Method != "GET" {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	w_products := Products{}

	// Get list of products from DB
	db.Find(&w_products.Products)

	// Transform to JSON format and send to a client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(w_products)
}

/**
* @api {get} /product/list/:id Get Product's Info
 * @apiName ListProduct
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
 * @apiError InvalidRequest The request is invalid.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 400 Bad Request
 *     {
 *       "error": "InvalidRequest"
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
 * @apiError (Server Error 5xx) InternalServerError There was a server-side error while processing the request.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *       "error": "InternalServerError"
 *     }
*/

func ListProduct(w http.ResponseWriter, r *http.Request) {
	// Check that method is allowed
	if r.Method != "GET" {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	// Query the database to get the product with the given ID
	fmt.Println("# Query from table products")
	product := Product{}
	result := db.First(&product, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(w, "ProductNotFound", http.StatusNotFound)
		return
	} else if result.Error != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	// Transform to JSON format and send to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
}

/**
* @api {post} /product/update/:id Update Product's Info
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
 * @apiError InvalidRequest The request is invalid.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 400 Bad Request
 *     {
 *       "error": "InvalidRequest"
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
 * @apiError (Server Error 5xx) InternalServerError There was a server-side error while processing the request.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 Internal Server Error
 *     {
 *       "error": "InternalServerError"
 *     }
*/

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Check that method is allowed
	if r.Method != "POST" {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	// Parse the product ID from the URL parameter
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	// Parse the product information from the request body
	var updatedProduct Product
	err = json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, "InvalidRequest", http.StatusBadRequest)
		return
	}

	// Update the product in the database
	fmt.Println("# Update product in table products")
	result := db.Model(&Product{}).Where("id = ?", id).Updates(updatedProduct)
	if result.RowsAffected == 0 {
		http.Error(w, "ProductNotFound", http.StatusNotFound)
		return
	}

	// Query the database to get the updated product information
	fmt.Println("# Query from table products")
	product := Product{}
	if err := db.First(&product, id).Error; err != nil {
		fmt.Println("Failed in db.First to DB")
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	// Return the updated product information to the client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{"product": product})
	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}
}
