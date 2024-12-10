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
    args := flag.CommandLine.Args()

    fmt.Printf("q path: %v\n", *q_path);
    fmt.Printf("db path: %v\n", *db_path);

	parser := newSqlParser();
	entries, err := os.ReadDir(*q_path);
	for _, entry := range entries {
        path := fmt.Sprintf("%s/%s", *q_path, entry.Name());
		err := parser.AddFromFile(path);
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

    if len(args) == 0{
        fmt.Println("no args!");
        os.Exit(-1);
    }
    switch args[0] {
        case "start server":
            server(database)
            break
        case "seed":
            if len(args) != 2{
                fmt.Println("missing seed path");
                os.Exit(-1);
            }
            seedDatabase(args[1], database);
            break;
        case "list_products":
            ps, _ := database.ListProducts();
            for _, p := range ps{
                fmt.Printf("%v\n", p);
            }

        case "query":
            if len(args) != 2{
                fmt.Println("missing query");
                os.Exit(-1);
            }
            rows, err := database.db.Query(args[1]);
            defer rows.Close();
            if err != nil{
                fmt.Printf("Error: %v\n", err);
                os.Exit(-1);
            }

            tys, err := rows.ColumnTypes()
            if err != nil{
                fmt.Printf("Error checking col types: %v\n", err);
                os.Exit(-1);
            }
            vals := []any {};
            for _, ty := range tys{
                switch ty.DatabaseTypeName(){
                    case "VARCHAR(256)":
                        fmt.Println("VARCHAR(256)");
                        vals = append(vals, "[STRING]");
                    break;
                    case "INTEGER":
                        fmt.Println("INTEGER");
                        vals = append(vals, 0)
                    break;
                    case "BOOL":
                        fmt.Println("BOOL");
                        vals = append(vals, false)
                    break;
                    default:
                        fmt.Printf("unknown ty: %v\n", ty.DatabaseTypeName());

                }
            }
            valsr := []any{}
            for _, val := range vals{
                valsr = append(valsr, &val);
            }
            for rows.Next(){
                err := rows.Scan(valsr...);
                if err != nil{
                    fmt.Printf("error while scanning: %f\n", err);
                    os.Exit(-1);
                }
                for _, val := range valsr{
                    fmt.Printf("%v\t", val)
                }
                fmt.Println();

            }
            break;
        case "exec":
            if len(args) != 2{
                fmt.Println("missing exec");
                os.Exit(-1);
            }
            database.db.Exec(args[1])
            break;
    }
}

func server(database *Database) {
	e := echo.New()

	e.Use(middleware.Logger());

	// non session api stuff
	bind_apis(e, database);


	assertHanlder := http.FileServer(http.FS(os.DirFS("static/")))
	e.GET("/*", echo.WrapHandler(http.StripPrefix("/", assertHanlder)))

	e.Logger.Fatal(e.Start(":8080"))
}
