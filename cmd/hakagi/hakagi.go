package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/syucream/hakagi/src/database"
	"github.com/syucream/hakagi/src/formatter"
	"github.com/syucream/hakagi/src/guess"
)

var ruleToGuesser = map[string]guess.GuessOption{
	"primarykey":     guess.GuessByPrimaryKey(),
	"tableandcolumn": guess.GuessByTableAndColumn(),
}

func main() {
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database password")
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.Int("dbport", 3306, "database port")

	targets := flag.String("targets", "", "analysing target databases(comma-separated)")
	rules := flag.String("rules", "primarykey,tableandcolumn", "analysing rules(comma-separated)")

	flag.Parse()

	db, err := database.ConnectDatabase(*dbUser, *dbPass, *dbHost, *dbPort)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	targetSlice := strings.Split(*targets, ",")
	schemas, err := database.FetchSchemas(db, targetSlice)
	if err != nil {
		log.Fatalf("Failed to fetch schemas : %v", err)
	}
	primaryKeys, err := database.FetchPrimaryKeys(db, targetSlice)
	if err != nil {
		log.Fatalf("Failed to fetch primary keys : %v", err)
	}

	var guessers []guess.GuessOption
	for _, rule := range strings.Split(*rules, ",") {
		if guesser, ok := ruleToGuesser[rule]; ok {
			guessers = append(guessers, guesser)
		}
	}
	constraints := guess.GuessConstraints(schemas, primaryKeys, guessers...)

	alterQuery := formatter.FormatSql(constraints)
	fmt.Println(alterQuery)
}
