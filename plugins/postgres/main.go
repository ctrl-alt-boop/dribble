package main

import (
	_ "database/sql"

	_ "github.com/lib/pq"
)

var Loaded bool = false

func main() {
	//
}

func init() {
	//
	Loaded = true
}
