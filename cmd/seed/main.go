package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"

	"exchange/internal/db"
)

type Pair struct {
	Ticker string
	Name   string
}

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432/testdb?sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}
	queries := db.New(database)

	log.Print("seeding database...")

	// Example seeding assets
	assets := []Pair{
		{Ticker: "BTC", Name: "Bitcoin"},
		{Ticker: "ETH", Name: "Ethereum"},
		{Ticker: "LINK", Name: "Chainlink"},
		{Ticker: "MADT", Name: "Tokenized MAD"},
	}

	ctx := context.Background()

	for _, asset := range assets {
		_, err := queries.CreateAsset(ctx, db.CreateAssetParams{
			AssetName: asset.Name,
			Ticker:    asset.Ticker,
		})
		if err != nil {
			log.Printf("error seeding asset %s: %v", asset.Ticker, err)
		} else {
			log.Printf("seeded asset: %s", asset.Ticker)
		}
	}

	queries.CreateUser(ctx, db.CreateUserParams{
		FirstName: "Imad",
		LastName:  "Archid",
		Dob:       time.Now(),
		Balance:   10000,
		Email:     "imad@exchange.co",
	})

	log.Println("done.")
}
