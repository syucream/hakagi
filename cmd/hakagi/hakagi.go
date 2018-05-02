package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	databaseName = "information_schema"
	queryBase    = `SELECT 
  TABLE_NAME, COLUMN_NAME, DATA_TYPE
FROM
  COLUMNS
WHERE
  TABLE_SCHEMA IN (%s);
`
)

type schema struct {
	table    string
	column   string
	dataType string
}

func fetchSchema(db *sql.DB, targets string) (schema, error) {
	targetsConditions := "'" + strings.Join(strings.Split(targets, ","), "','") + "'"

	rows, err := db.Query(queryBase, targetsConditions)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

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

	dbUri := fmt.Sprinft("%s:%s@tcp(%s:%d)/%s", *dbUser, *dbPass, *dbhost, *dbPort, databaseName)
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
	}

	schemas, err := fetchSchema(db, *targets)
	if err != nil {
		log.Fatalf("Failed to fetch schemas : %v", err)
	}
}
