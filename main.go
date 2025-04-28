package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"exchange/internal/api/handler"
	"exchange/internal/api/router"
	"exchange/internal/db"
)

func main() {

	connStr := "postgresql://postgres:postgres@localhost:5432/testdb?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	queries := db.New(database)

	h := &handler.Handler{Queries: queries}
	r := router.NewRouter(h)

	http.ListenAndServe(":3000", r)
}
