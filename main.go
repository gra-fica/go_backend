package main

import (
	"./parser"
)

/*
import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"database/sql"
)

type id_t = uint64

type Product struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
	Price uint64 `json:"price"`
}

type Alias struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
	ProductID id_t `json:"product_id"`
}

type Category struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
}

type Sale struct {
	ID id_t `json:"id"`
	ProductID id_t `json:"product_id"`
	Quantity int `json:"quantity"`
	Price uint64 `json:"price"`
}

type Ticket struct {
	ID id_t `json:"id"`
	SaleID []id_t `json:"sale_id"`
	Total uint64 `json:"total"`
}

type Database struct {
	db *sql.DB
}

func (*Database)newProduct(_name string, _price uint64) error {
	// Add product to database
	// Return error if product already exists

	return nil;
}

func (*Database)newAlias(_name string, _product_id id_t) error {
	// Add alias to database
	// Return error if alias already exists

	return nil;
}

func (*Database)newCategory(_name string) error {
	// Add category to database
	// Return error if category already exists

	return nil;
}

func (*Database)newSale(_product_id id_t, _quantity int, _price uint64) error {
	// Add sale to database
	// Return error if sale already exists

	return nil;
}

func (*Database)newTicket(_sale_id []id_t, _total uint64) error {
	// Add ticket to database
	// Return error if ticket already exists

	return nil;
}

func initDatabase() *Database {
	database := Database{};
	db, err := sql.Open("sqlite3", "./database.db");
	if err != nil {
		panic(err);
	}

	database.db = db;

	return &database;
}
*/

func main(){
	/*
	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/api/v1/product/list", func (c echo.Context) error {
		return c.String(200, "List of products")
	});
	*/

}

