package main

import (
	"flag"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"os"

	_ "github.com/mattn/go-sqlite3" // yo go wtf

	"fmt"
)

func main(){
    q_path := flag.String("q", "./sql/", "where the queries are at");
    db_path := flag.String("db", "database.sql", "where the database is");

    flag.Parse()

    fmt.Printf("q path: %v\n", q_path);
    fmt.Printf("db path: %v\n", db_path);

	parser := newSqlParser();
	entries, err := os.ReadDir(*q_path);
	for _, entry := range entries {
		err := parser.AddFromFile(entry.Name());
		if err != nil {
            fmt.Printf("Error parsing sql: %v", err);
			return
		}
	}

	database, err := initDatabase(parser, *db_path);
	if err != nil {
        fmt.Println("Error initing database");
		return
	}

	e := echo.New()

	e.Use(middleware.Logger());

	// non session api stuff
	bind_apis(e, database);


	assertHanlder := http.FileServer(http.FS(os.DirFS("static/")))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assertHanlder)))

	e.Logger.Fatal(e.Start(":8080"))
}
