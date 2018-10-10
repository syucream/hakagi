package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	databaseName     = "information_schema"
	indexesQueryBase = `SELECT
  STATISTICS.TABLE_NAME, STATISTICS.INDEX_NAME, STATISTICS.COLUMN_NAME, COLUMNS.COLUMN_TYPE
FROM
  STATISTICS
INNER JOIN
	COLUMNS
ON
  STATISTICS.TABLE_NAME = COLUMNS.TABLE_NAME AND
	STATISTICS.COLUMN_NAME = COLUMNS.COLUMN_NAME AND
    STATISTICS.TABLE_SCHEMA = COLUMNS.TABLE_SCHEMA
WHERE
  STATISTICS.INDEX_NAME != 'PRIMARY' AND
  STATISTICS.TABLE_SCHEMA IN (?);
`
	primaryConstraintQueryBase = `SELECT
  KEY_COLUMN_USAGE.TABLE_NAME, KEY_COLUMN_USAGE.COLUMN_NAME, COLUMNS.COLUMN_TYPE
FROM
  KEY_COLUMN_USAGE
INNER JOIN
	COLUMNS
ON
  KEY_COLUMN_USAGE.TABLE_NAME = COLUMNS.TABLE_NAME AND
	KEY_COLUMN_USAGE.COLUMN_NAME = COLUMNS.COLUMN_NAME AND
    KEY_COLUMN_USAGE.CONSTRAINT_SCHEMA = COLUMNS.TABLE_SCHEMA
WHERE
  KEY_COLUMN_USAGE.CONSTRAINT_NAME = 'PRIMARY' AND
  KEY_COLUMN_USAGE.CONSTRAINT_SCHEMA IN (?);
`
)

type Column struct {
	Name string
	Type string
}

// table name -> index name -> columns
type Indexes map[string]map[string][]Column

// table name -> columns
type PrimaryKeys map[string][]Column

func ConnectDatabase(user string, pass string, host string, port int) (*sql.DB, error) {
	dbUri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, databaseName)

	return sql.Open("mysql", dbUri)
}

func FetchSchemas(db *sql.DB, targets []string) (Indexes, error) {
	query, args, err := sqlx.In(indexesQueryBase, targets)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexes := make(Indexes)
	for rows.Next() {
		var tableName, indexName, columnName, columnType string
		if err := rows.Scan(&tableName, &indexName, &columnName, &columnType); err != nil {
			return nil, err
		}

		if ct, tableOk := indexes[tableName]; tableOk {
			if ci, indexOk := ct[indexName]; indexOk {
				ct[indexName] = append(ci, Column{columnName, columnType})
			} else {
				ct[indexName] = []Column{Column{columnName, columnType}}
			}
		} else {
			indexes[tableName] = make(map[string][]Column)
			indexes[tableName][indexName] = []Column{Column{columnName, columnType}}
		}
	}

	return indexes, nil
}

func FetchPrimaryKeys(db *sql.DB, targets []string) (PrimaryKeys, error) {
	query, args, err := sqlx.In(primaryConstraintQueryBase, targets)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	primaryKeys := make(PrimaryKeys)
	for rows.Next() {
		var tableName, columnName, columnType string
		if err := rows.Scan(&tableName, &columnName, &columnType); err != nil {
			return nil, err
		}

		if current, ok := primaryKeys[tableName]; ok {
			primaryKeys[tableName] = append(current, Column{columnName, columnType})
		} else {
			primaryKeys[tableName] = []Column{Column{columnName, columnType}}
		}
	}

	return primaryKeys, nil
}
