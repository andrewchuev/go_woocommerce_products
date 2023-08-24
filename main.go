package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

var db *sql.DB

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Dbname,
	)

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
	baseQuery := `
	SELECT
	    p.ID,
	    p.post_title AS name,
	    MAX(CASE WHEN m.meta_key = '_price' THEN m.meta_value ELSE NULL END) AS price,
	    p.post_content AS description,
	    GROUP_CONCAT(DISTINCT t.name ORDER BY t.name ASC) AS categories
	FROM
	    wp_posts p
	LEFT JOIN
	    wp_postmeta m ON p.ID = m.post_id AND m.meta_key = '_price'
	LEFT JOIN
	    wp_term_relationships tr ON p.ID = tr.object_id
	LEFT JOIN
	    wp_term_taxonomy tt ON tr.term_taxonomy_id = tt.term_taxonomy_id AND tt.taxonomy = 'product_cat'
	LEFT JOIN
	    wp_terms t ON tt.term_id = t.term_id
	WHERE
	    p.post_type = 'product'
	`

	// Фильтрация
	params := r.URL.Query()
	if category, ok := params["category"]; ok && len(category) > 0 {
		baseQuery += fmt.Sprintf(" AND t.name = '%s'", category[0])
	}
	if minPrice, ok := params["min_price"]; ok && len(minPrice) > 0 {
		baseQuery += fmt.Sprintf(" AND m.meta_value >= '%s'", minPrice[0])
	}
	if maxPrice, ok := params["max_price"]; ok && len(maxPrice) > 0 {
		baseQuery += fmt.Sprintf(" AND m.meta_value <= '%s'", maxPrice[0])
	}

	// Группировка
	baseQuery += " GROUP BY p.ID"

	// Сортировка
	if sort, ok := params["sort"]; ok && len(sort) > 0 {
		order := "asc"
		if o, ok := params["order"]; ok && len(o) > 0 && (o[0] == "desc" || o[0] == "asc") {
			order = o[0]
		}
		switch sort[0] {
		case "price":
			baseQuery += fmt.Sprintf(" ORDER BY price %s", order)
		case "name":
			baseQuery += fmt.Sprintf(" ORDER BY name %s", order)
		}
	}

	rows, err := db.Query(baseQuery)
	if err != nil {
		http.Error(w, "Unable to query database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Category)
		if err != nil {
			http.Error(w, "Unable to scan row", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Dbname   string `json:"dbname"`
	} `json:"database"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
