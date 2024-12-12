package main

import (
	"fmt"
)

func query(d *Database, query string) (err error) {
	rows, err := d.Query(query);
    fmt.Printf("rows: %v", rows);
	return
}

func exec(d *Database, query string) (err error) {
	result, err := d.Execute(query);
    fmt.Printf("result: %v", result);
	return
}

func cmd(args []string) (err error) {
	return
}
