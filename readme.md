# WooCommerce Products Microservice

## This microservice provides an API to retrieve a list of products from a WordPress database where the WooCommerce plugin is installed.


### Getting Started

#### Prerequisites

Go installed (version 1.16 or higher)
A WordPress database with the WooCommerce plugin installed

#### Installation

Clone the repository:

    git clone https://github.com/andrewchuev/go_woocommerce_products
    cd woocommerce-microservice

Install the required dependencies:

    go get -u github.com/go-sql-driver/mysql

Configure the database connection parameters in the config.json file:

    {
        "database": {
        "username": "your_username",
        "password": "your_password",
        "host": "127.0.0.1",
        "port": "3306",
        "dbname": "wordpress"
        }
    }

Replace your_username and your_password with your actual values.


#### Running

    go run main.go

After starting, the microservice will be available at http://localhost:8080.
To retrieve the list of products, send a GET request to http://localhost:8080/products.


### License

This project is licensed under the MIT License - see the LICENSE.md file for details.