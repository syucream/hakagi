package guess

import (
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/syucream/hakagi/src/constraint"
	"github.com/syucream/hakagi/src/database"
)

const (
	idColumn           = "id"
	targetColumnSuffix = "_id"
)

type GuessOption func(database.Column, string, database.Column) bool

func isAcceptableAsPrimaryKey(columnType, primaryKeyType string) bool {
	colIsOk := strings.Index(columnType, "int") != -1
	pkIsOk := strings.Index(primaryKeyType, "int") != -1
	return colIsOk && pkIsOk && columnType == primaryKeyType
}

// Recongnize a column thats same name of other table's primary key is a foreign key
// This base idea refers to SchemaSpy DbAnalyzer:
//   https://github.com/schemaspy/schemaspy/blob/master/src/main/java/org/schemaspy/DbAnalyzer.java
func GuessByPrimaryKey() GuessOption {
	return func(i database.Column, table string, pk database.Column) bool {
		return isAcceptableAsPrimaryKey(i.Type, pk.Type) && i.Name == pk.Name && pk.Name != idColumn
	}
}

func GuessByTableAndColumn() GuessOption {
	return func(i database.Column, table string, pk database.Column) bool {
		if !isAcceptableAsPrimaryKey(i.Type, pk.Type) {
			return false
		}

		cLen := len(i.Name)
		tLen := len(targetColumnSuffix)
		if !(cLen >= tLen && i.Name[cLen-tLen:] == targetColumnSuffix) {
			return false
		}

		return inflection.Plural(i.Name[:cLen-tLen]) == table && pk.Name == idColumn
	}
}

// GuessConstraints guesses foreign key constraints from primary keys and indexes.
// NOTE composite primary keys are not supported.
func GuessConstraints(indexes database.Indexes, primaryKeys database.PrimaryKeys, guessOptions ...GuessOption) []constraint.Constraint {
	var constraints []constraint.Constraint

	for indexTable, indexMaps := range indexes {
		for _, indexCols := range indexMaps {
			for pkTable, pk := range primaryKeys {
				if indexTable != pkTable && len(indexCols) == 1 && len(pk) == 1 {
					singleIndex := indexCols[0]
					singlePk := pk[0]

					for _, guesser := range guessOptions {
						if guesser(singleIndex, pkTable, singlePk) {
							constraints = append(constraints, constraint.Constraint{indexTable, singleIndex.Name, pkTable, singlePk.Name})
						}
					}
				}
			}
		}
	}

	return constraints
}
