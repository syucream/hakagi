package guess

import (
	"github.com/jinzhu/inflection"
	"github.com/syucream/hakagi/src/constraint"
	"github.com/syucream/hakagi/src/database"
)

const (
	idColumn           = "id"
	targetColumnSuffix = "_id"
)

var primaryKeyTypes = map[string]bool{
	"tinyint":   true,
	"smallint":  true,
	"meriumint": true,
	"int":       true,
	"bigint":    true,
}

type GuessOption func(database.Schema, string, database.PrimaryKey) bool

func isAcceptableAsPrimaryKey(columnType, primaryKeyType string) bool {
	_, colIsOk := primaryKeyTypes[columnType]
	_, pkIsOk := primaryKeyTypes[primaryKeyType]
	return colIsOk && pkIsOk && columnType == primaryKeyType
}

// Recongnize a column thats same name of other table's primary key is a foreign key
// This base idea refers to SchemaSpy DbAnalyzer:
//   https://github.com/schemaspy/schemaspy/blob/master/src/main/java/org/schemaspy/DbAnalyzer.java
func GuessByPrimaryKey() GuessOption {
	return func(s database.Schema, table string, pk database.PrimaryKey) bool {
		return isAcceptableAsPrimaryKey(s.DataType, pk.DataType) && s.Column == pk.Column && pk.Column != idColumn
	}
}

func GuessByTableAndColumn() GuessOption {
	return func(s database.Schema, table string, pk database.PrimaryKey) bool {
		if !isAcceptableAsPrimaryKey(s.DataType, pk.DataType) {
			return false
		}

		cLen := len(s.Column)
		tLen := len(targetColumnSuffix)
		if !(cLen >= tLen && s.Column[cLen-tLen:] == targetColumnSuffix) {
			return false
		}

		return inflection.Plural(s.Column[:cLen-tLen]) == table && pk.Column == idColumn
	}
}

func GuessConstraints(schemas []database.Schema, primaryKeys database.PrimaryKeys, guessOptions ...GuessOption) []constraint.Constraint {
	var constraints []constraint.Constraint

	for _, s := range schemas {
		for table, pk := range primaryKeys {
			// NOTE composite primary keys are not supported
			if s.Table != table && len(pk) == 1 {
				singlePk := pk[0]
				for _, guesser := range guessOptions {
					if guesser(s, table, singlePk) {
						constraints = append(constraints, constraint.Constraint{s.Table, s.Column, table, singlePk.Column})
					}
				}
			}
		}
	}

	return constraints
}
