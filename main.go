package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"


	_ "github.com/mattn/go-sqlite3" // yo go wtf

	"fmt"
)

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
	e.Use(middleware.Logger());

	// non session api stuff
	bind_apis(e, database);
	e.Logger.Fatal(e.Start(":8080"))
}
