package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

var db *sql.DB

func main() {
	var err error
	var dsn = "homestead:secret@tcp(192.168.88.56:3306)/urbanlab2"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", getProducts)

	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT ID, post_title, meta_value FROM wp_posts INNER JOIN wp_postmeta ON wp_posts.ID = wp_postmeta.post_id WHERE post_type = 'product' AND meta_key = '_price'")
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			http.Error(w, "Unable to scan row", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
