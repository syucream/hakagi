package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	databaseName = "information_schema"
	queryBase    = `SELECT 
  TABLE_NAME, COLUMN_NAME, DATA_TYPE
FROM
  COLUMNS
WHERE
  TABLE_SCHEMA IN (?);
`
)

type schema struct {
	table    string
	column   string
	dataType string
}

func fetchSchemas(db *sql.DB, targets []string) ([]schema, error) {
	query, args, err := sqlx.In(queryBase, targets)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println(rows)

	var schemas []schema
	for rows.Next() {
		var tableName, columnName, dataType string
		if err := rows.Scan(&tableName, &columnName, &dataType); err != nil {
			return nil, err
		}
		schemas = append(schemas, schema{tableName, columnName, dataType})
	}

	return schemas, nil
}

func main() {
	dbUser := flag.String("dbuser", "", "database user")
	dbPass := flag.String("dbpass", "", "database password")
	dbHost := flag.String("dbhost", "localhost", "database host")
	dbPort := flag.Int("dbport", 3306, "database port")

	targets := flag.String("targets", "", "analysing targets(comma-separated)")

	flag.Parse()

	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *dbUser, *dbPass, *dbHost, *dbPort, databaseName)
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	schemas, err := fetchSchemas(db, strings.Split(*targets, ","))
	if err != nil {
		log.Fatalf("Failed to fetch schemas : %v", err)
	}
	fmt.Println(schemas)
}
