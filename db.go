package main

import (
	"encoding/json"

	"crypto/sha256"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // yo go wtf

	"database/sql"
	"fmt"
)

type id_t = uint64

type User struct{
	ID       id_t   `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type UserInfo struct{
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Product struct {
	ID    id_t   `json:"id"`
	Name  string `json:"name"`
	Price uint64 `json:"price"`
	Cost  uint64 `json:"cost"`
	Count uint64 `json:"count"`
}

type ParcialProduct struct {
	ID    id_t   `json:"id"`
	Name  string `json:"name"`
	Price uint64 `json:"price"`
}

type ProductBuffer struct {
	Name    string `json:"name"`
	Aliases []string `json:"aliases"`
	Price   uint64 `json:"price"`
	Cost    uint64 `json:"cost"`
	Count   uint64 `json:"count"`
} 

type Alias struct {
	ID        id_t `json:"id"`
	Name      string `json:"name"`
	ProductID id_t `json:"product_id"`
}

type Category struct {
	ID   id_t `json:"id"`
	Name string `json:"name"`
}

type Sale struct {
	ID        id_t `json:"id"`
	ProductID id_t `json:"product_id"`
	Quantity  int `json:"quantity"`
	Price     uint64 `json:"price"`
}

type Ticket struct {
	ID     id_t `json:"id"`
	SaleID []id_t `json:"sale_id"`
	Total  uint64 `json:"total"`
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
		"USER",
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
	rows, e := d.Query("LIST-PRODUCTS-PARCIAL");
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

func (p *Product) GetStringFuzzy() *string {
	return &p.Name;
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

func (d *Database) CreateTicket() (t Ticket, err error){
	response , err := d.Execute("CREATE-TICKET");
	if err != nil {
		return
	}
	id, err := response.LastInsertId();
	if err != nil{
		return
	}
	q, err := d.Query("GET-TICKET-ID", id);
	q.Scan(&t.ID, &t.Total, &t.SaleID);

	return
}

// badly made login and signin functionsfunctions
func (d* Database) SignIn(name string, password string) (res sql.Result, err error){
	row, err := d.Query("FIND-USERNAME", name);
	fmt.Printf("using adding user %s: %s\n", name, password);

	// if there is a query then the name exists
	if row.Next() {
		name = "";
		row.Scan(&name);
		fmt.Printf("found name: %s\n", name);
		err = fmt.Errorf("there already exists this name");
		fmt.Println("there already exists this name!!");
		return
	}

	hash := sha256.New();
	_, err = hash.Write([]byte(password))
	if err != nil {
		fmt.Println("could not hash the password??");
		return
	}

	hashed_password := fmt.Sprintf("%x", hash.Sum(nil));
	res, err = d.Execute("ADD-USER", name, "default", hashed_password);
	if err != nil{
		fmt.Printf("failed to add user: %v\n", err);
	}	
	return;
}

func (d* Database) LogIn(name string, password string) (ok bool, err error){
	hash := sha256.New();
	_, err = hash.Write([]byte(password))
	row, err := d.Query("FIND-USER", name);
	if err != nil{
		fmt.Println("couldn not find user");
		return;
	}
	hashed_password := fmt.Sprintf("%x", hash.Sum(nil));
	user := User{};
	row.Next()
	err = row.Scan(&user.ID, &user.Name, &user.Role, &user.Password);

	if err != nil{
		return
	}

	ok = hashed_password == user.Password;
	return;
}
