package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	databaseName    = "information_schema"
	schemaQueryBase = `SELECT 
  TABLE_NAME, COLUMN_NAME, DATA_TYPE
FROM
  COLUMNS
WHERE
  TABLE_SCHEMA IN (?);
`
	primaryConstraintQueryBase = `SELECT 
	TABLE_NAME, COLUMN_NAME
FROM
  KEY_COLUMN_USAGE
WHERE
  CONSTRAINT_NAME = 'PRIMARY' AND
  CONSTRAINT_SCHEMA IN (?);
`
)

type Schema struct {
	Table    string
	Column   string
	DataType string
}

type PrimaryKey struct {
	Table  string
	Column string
}

func ConnectDatabase(user string, pass string, host string, port int) (*sql.DB, error) {
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, databaseName)

	return sql.Open("mysql", dbUri)
}

func FetchSchemas(db *sql.DB, targets []string) ([]Schema, error) {
	query, args, err := sqlx.In(schemaQueryBase, targets)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []Schema
	for rows.Next() {
		var tableName, columnName, dataType string
		if err := rows.Scan(&tableName, &columnName, &dataType); err != nil {
			return nil, err
		}
		schemas = append(schemas, Schema{tableName, columnName, dataType})
	}

	return schemas, nil
}

func FetchPrimaryKeys(db *sql.DB, targets []string) ([]PrimaryKey, error) {
	query, args, err := sqlx.In(primaryConstraintQueryBase, targets)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKeys []PrimaryKey
	for rows.Next() {
		var tableName, columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, PrimaryKey{tableName, columnName})
	}

	return primaryKeys, nil
}
