package database

import (
	ent "blog/ent"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Vars map[string]string

func connection() *ent.Client {
	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Vars["PG_HOST"], Vars["PG_PORT"], Vars["PG_USER"], Vars["PG_DB"], Vars["PG_PASS"])

	client, err := ent.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	return client
}
