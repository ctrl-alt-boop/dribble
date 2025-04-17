package main

import (
	_ "database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var Loaded bool = false

func main() {
	//
}

func init() {
	//
	Loaded = true
}
