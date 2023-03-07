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

/**
TASK:
При помощи инструмента apiDoc написать документацию к API абстрактного интернет-магазина.
Опишите минимум 3 метода HTTP для следующих операций:

- Просмотр всех товаров в магазине.
GET /products - Выполнено
- Вывод детальной информации о конкретном товаре.
GET /product/:id - Выполнено
- Изменение информации о товаре (доступно только администратору).
POST /product/:id - Выполнено
- Добавление товара в корзину.
POST /checkout - Не выполнено
*/

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
 *           "id": 1,
 *           "name": "Product 1",
 *           "description": "Product 1 description",
 *           "price": 19.99
 *         },
 *         {
 *           "id": 2,
 *           "name": "Product 2",
 *           "description": "Product 2 description",
 *           "price": 29.99
 *         },
 *         {
 *           "id": 3,
 *           "name": "Product 3",
 *           "description": "Product 3 description",
 *           "price": 39.99
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

func getProducts(w http.ResponseWriter, r *http.Request) {
	w_products := Products{}

	// Establish connection to container DB
	db, err := sql.Open("postgres", "host=postgresql user=user-api password=qwe123 dbname=store_api sslmode=disable")
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

/**
* @api {get} /product/:id Get Product's Info
 * @apiName GetProductInfo
 * @apiGroup Products
 * @apiPermission user
 *
 * @apiParam {Number} id Products unique ID.
 *
 * @apiSuccess {Number} id Products unique ID.
 * @apiSuccess {String} name Name of the Product.
 * @apiSuccess {String} description Description of the Product.
 * @apiSuccess {Number} price Price of the Product.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "product": {
 *         "id": 123,
 *         "name": "Product name",
 *         "description": "Product description",
 *         "price": 99.99
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

/**
* @api {post} /product/:id Update Product's Info
 * @apiName UpdateProductInfo
 * @apiGroup Products
 * @apiPermission admin
 *
 * @apiParam {Number} id Products unique ID.
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
 *         "id": 123,
 *         "name": "New name",
 *         "description": "New description",
 *         "price": 100
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
