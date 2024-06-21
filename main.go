package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // yo go wtf

	"database/sql"
	"fmt"
)

type id_t = uint64

type Product struct {
	ID id_t `json:"id"`
	Name string `json:"name"`
	Price uint64 `json:"price"`
}

type ProductBuffer struct {
	Name string `json:"name"`
	Price uint64 `json:"price"`
	Aliases []string `json:"aliases"`
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
	db      *sql.DB
	parser  *SqlParser
}

func (d *Database) Execute(query string, params ...any) (r sql.Result, err error) {
	if d.parser == nil{
		err = fmt.Errorf("No parser found");
		return
	}

	q, ok := d.parser.formats[query]
	if !ok {
		err = fmt.Errorf("Query not found");
		return
	}
	r, err = d.db.Exec(q, params...);
	if err != nil {
		err = fmt.Errorf("could not query [%s] error: [%v]", query, err);
		return
	}
	return
}

func (d *Database) Query(query string, params ...any) (r *sql.Rows, err error) {
	if d.parser == nil{
		err = fmt.Errorf("No parser found");
		return
	}

	q, ok := d.parser.formats[query]
	if !ok {
		err = fmt.Errorf("Query not found");
		return
	}
	raw, err := d.db.Query(q, params...);
	if err != nil {
		err = fmt.Errorf("could not query [%s] error: [%v]", query, err);
		return
	}

	if raw == nil {
		err = fmt.Errorf("empty query [%s]", query);
		return 
	}

	r = raw;
	return
}

func (d *Database)NewProduct(_name string, _price uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-PRODUCT", _name, _price);
	return
}

func (d *Database)newAlias(_name string, _product_id id_t) (r sql.Result, err error) {
	r, err = d.Execute("ADD-ALIAS", _name, _product_id);
	return
}

func (d *Database)newCategory(_name string) (r sql.Result, err error) {
	r, err = d.Execute("ADD-CATEGORY", _name);
	return
}

func (d *Database)newSale(_product_id id_t, _quantity int, _price uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-SALE", _product_id, _quantity, _price);
	return
}

func (d *Database)newTicket(_sale_id []id_t, _total uint64) (r sql.Result, err error) {
	r, err = d.Execute("ADD-TICKET", _sale_id, _total);
	return
}

func initDatabase(p *SqlParser) (database *Database, err error){
	db, err := sql.Open("sqlite3", "./database.db");

	database = &Database{
		parser: p,
	};
	if err != nil {
		panic(err);
	}

	database.db = db;

	// Create tables
	tables := []string {
		"PRODUCT",
		"ALIAS",
		"CATEGORY",
		"SALE",
		"TICKET",
	};
	for _, table := range tables {
		_, err = database.Execute("CREATE-" + table);
		if err != nil {
			fmt.Printf("COULD NOT CREATE TABLE: %s %v\n", table, err);
			return
		}
	}

	return
}

func seedDatabase(database *Database) (err error) {
	path := "./.ignore/seed.csv";
	data, err := os.ReadFile(path);
	if err != nil {
		err = fmt.Errorf("could not open file %s", path);
		return
	}

	data_str := string(data);
	lines := strings.Split(data_str, "\n");

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Split(line, ";");
		if len(parts) == 0 {
			err = fmt.Errorf("could not parse line %s", line);
			return
		}

		name := parts[0];
		price, err := strconv.ParseUint(parts[1], 10, 64);
		if err != nil {
			err = fmt.Errorf("could not parse price %s", parts[1]);
			return err
		}

		_, err = database.NewProduct(name, price);
	}
	return
}

func (d *Database) ListProductsWherePrice(price int) (p []Product, e error) {
	rows, e := d.Query("LIST-PRODUCTS-WHERE-PRICE-IS", price);
	if e != nil {
		return	
	}

	for rows.Next() {
		var product Product;
		e = rows.Scan(&product.ID, &product.Name, &product.Price);
		if e != nil {
			return
		}
		p = append(p, product);
	}

	return
}

func (d* Database) GetProductFromName(name string) (p Product, err error){
	rows, err := d.Query("GET-PRODUCT-NAME", name);
	rows.Scan(&p.ID, &p.Name, &p.Price);

	return
}

func (d *Database) ListProducts() (p []Product, e error) {
	rows, e := d.Query("LIST-PRODUCTS");
	if e != nil {
		return	
	}

	for rows.Next() {
		var product Product;
		e = rows.Scan(&product.ID, &product.Name, &product.Price);
		if e != nil {
			return
		}
		p = append(p, product);
	}

	return
}

func (d* Database) AddProduct(p ProductBuffer) (err error){
	if len(p.Name) == 0{
		err = &json.UnsupportedValueError{}
		return
	}
	_, err = d.Execute("ADD-PRODUCT", p.Name, p.Price);
	if err != nil {
		return
	}
	// for _ = range p.Aliases { }
 
	return
}

func (d *Database) DeleteProduct(id id_t) (r sql.Result, err error) {
	r, err = d.Execute("DELETE-PRODUCT", id);
	return
}

func (d *Database) DeleteAlias(id id_t) (r sql.Result, err error) {
	r, err = d.Execute("DELETE-ALIAS", id);
	return
}

func (d *Database) DeleteCategory(id id_t) (r sql.Result, err error) {
	r, err = d.Execute("DELETE-CATEGORY", id);
	return
}

func (d *Database) DeleteSale(id id_t) (r sql.Result, err error) {
	r, err = d.Execute("DELETE-SALE", id);
	return
}

func (d *Database) DeleteTicket(id id_t) (r sql.Result, err error) {
	r, err = d.Execute("DELETE-TICKET", id);
	return
}

func (d *Database) GetProduct(id id_t) (p Product, e error) {
	rows, e := d.Query("GET-PRODUCT-ID");
	if e != nil {
		return	
	}

	for rows.Next() {
		e = rows.Scan(&p.ID, &p.Name, &p.Price);
		if e != nil {
			return
		}
	}

	return
}

type DatabaseLister interface {
	ListProducts(*Database) ([]Product, error)
}

type ListAllProducts struct {}

func (l *ListAllProducts) ListProducts(d *Database) ([]Product, error) {
	return d.ListProducts();
}

type ListProductsWherePrice struct {
	Price int
}

func (l *ListProductsWherePrice) ListProducts(d *Database) ([]Product, error) {
	return d.ListProductsWherePrice(l.Price);
}

type ProductMatch struct {
	Product Product `json:"product"`
	Score int       `json:"score"`
}

func (d *Database) SearchProduct(name string, lister DatabaseLister, matcher FuzzySearcher) (p []ProductMatch, e error) {
	prods, e := lister.ListProducts(d);
	if e != nil {
		return
	}

	bufp := []FuzzyObject{};
	for _, prod := range prods{
		bufp = append(bufp, &prod);
	}

	scores, e := matcher.Search(name, bufp);
	if e != nil {
		e = fmt.Errorf("could not search %v", e);
		return
	}
	for _, score := range scores {
		p = append(p, ProductMatch{*(score.Match.(*Product)), score.Score});
	}
	return
}

func (p *Product) GetStringFuzzy() *string {
	return &p.Name;
}

type Templates struct {
	ts *template.Template
}

func NewTemplates() *Templates{
	return &Templates{
		ts: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.ts.ExecuteTemplate(w, name, data)
}

func main(){
	parser := newSqlParser();
	files := []string {
		"./sql/add.sql",
		"./sql/create_table.sql",
		"./sql/list.sql",
		"./sql/delete.sql",
	};
	for _, file := range files {
		err := parser.AddFromFile(file);
		if err != nil {
			fmt.Println(err);
			return
		}
	}

	database, err := initDatabase(parser);
	if err != nil {
		return
	}

	e := echo.New()
	e.Renderer = NewTemplates();

	e.Use(middleware.Logger())


	e.GET("/api/v1/product/list", func (c echo.Context) error {
		products, err := database.ListProducts();
		if err != nil {
			return c.String(500, err.Error());
		}
		type ProductResponse struct {
			Products []Product `json:"products"`
		}
		jsonProds, err := json.Marshal(ProductResponse{products});
		return c.String(200, string(jsonProds))
	});

	e.GET("/api/v1/product/search/:name", func (c echo.Context) error {
		name := c.Param("name");

		products, err := database.SearchProduct(name, &ListAllProducts{}, &TokenFuzzy{});
		if err != nil {

			return c.String(500, err.Error());
		}
		type SearchProductResponse struct {
			Products []ProductMatch `json:"products"`
		}
		jsonProds, err := json.Marshal(SearchProductResponse{products});
		return c.String(200, string(jsonProds))
	});

	e.GET("/api/v1/product/search_exact/:name", func (c echo.Context) error {
		name := c.Param("name");
		product, err := database.GetProductFromName(name);
		if err != nil{
			return c.String(404, fmt.Sprintf("%s not found", name));
		}
		marshal, err := json.Marshal(product);
		if err != nil{
			return c.String(404, fmt.Sprintf("%s not found", name));
		}

		return c.String(200, string(marshal));
	})

	e.GET("/api/v1/product/search/:name/price/:price", func (c echo.Context) error {
		name := c.Param("name");
		price, err := strconv.Atoi(c.Param("price"));
		if err != nil {
			return c.String(400, fmt.Sprintf("price {%v} is not a number", c.Param("price")));
		}

		products, err := database.SearchProduct(name, &ListProductsWherePrice{price}, &TokenFuzzy{});
		type SearchProductResponse struct {
			Products []ProductMatch `json:"products"`
		}
		jsonProds, err := json.Marshal(SearchProductResponse{products});
		return c.String(200, string(jsonProds)) })


	e.DELETE("/api/v1/product/delete/:id", func (c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64);
		if err != nil {
			return c.String(400, fmt.Sprintf("id {%v} is not a number", c.Param("id")));
		}

		_, err = database.DeleteProduct(id);
		if err != nil {
			return c.String(500, err.Error());
		}
		return c.String(200, "ok")
	});

	e.POST("/api/v1/product/add", func (c echo.Context) error {

		payload := ProductBuffer {};
		err :=  (&echo.DefaultBinder{}).BindBody(c, &payload);

		if err != nil{
			return c.String(400, "Malfromated Product");
		}

		err = database.AddProduct(payload);
		if err != nil{
			return c.String(400, "failed to add product");
		}

		ans, err := json.Marshal(payload);
		if err != nil{
			return c.String(400, "Something very wrong just happend");
		}
		return c.String(200, string(ans));
	});


	e.GET("/", func(c echo.Context) error {
		return c.Redirect(300, "/index");
	})

	e.GET("/index", func(c echo.Context) error {
		return c.Render(200, "index", nil);
	})


	e.GET("/ticket", func(c echo.Context) error {
		return c.Render(200, "ticket", nil);
	})

	e.GET("/add_product", func(c echo.Context) error {
		return c.Render(200, "add_product", nil);
	})

	assertHanlder := http.FileServer(http.FS(os.DirFS("static/")))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assertHanlder)))

	e.GET("/api/v1/htmx/search/product", func(c echo.Context) error {
		name  := c.FormValue("name")
		fmt.Printf("name: %s\n", name)
		type ProductsMatch struct {
			Matches []ProductMatch
		}
		if len(name) <= 4 {
			fmt.Println("name under 4 letters")
			return c.Render(200, "search-result", ProductsMatch{});
		}
		
		prods, err := database.SearchProduct(name, &ListAllProducts{}, &TokenFuzzy{});
		for _, m := range prods{
			fmt.Printf("%d %s\n", m.Score, m.Product.Name);
		}
		matches := ProductsMatch {prods[:10]}
		if err != nil {
			return err
		}

		return c.Render(200, "search-result", matches)
	});

	e.Logger.Fatal(e.Start(":8080"))
}
