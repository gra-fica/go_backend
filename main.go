package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-jwt/v4"

	"os"

	_ "github.com/mattn/go-sqlite3" // yo go wtf

	"fmt"
)


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

	e.Use(middleware.Logger());
	e.Use(echojwt.JWT([]byte("secret")));

	// non session api stuff
	bind_apis(e, database);


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

	// index redirect stuff
	redirect := func(c echo.Context) error {
		return c.Redirect(300, "/index");
	};

	e.GET("", redirect);
	e.GET("/", redirect);

	// web page rendering
	e.GET("/index", func(c echo.Context) error {
		return c.Render(200, "index", nil);
	})

	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search", nil);
	})

	e.GET("/add_product", func(c echo.Context) error {
		return c.Render(200, "add_product", nil);
	})


	e.Logger.Fatal(e.Start(":8080"))
}
