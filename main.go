package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
		"./sql/user.sql",
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

	// non session api stuff
	bind_apis(e, database);


	assertHanlder := http.FileServer(http.FS(os.DirFS("static/")))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assertHanlder)))

	e.Logger.Fatal(e.Start(":8080"))
}
