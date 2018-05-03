package guess

import (
	"github.com/jinzhu/inflection"
	"github.com/syucream/hakagi/src/constraint"
	"github.com/syucream/hakagi/src/database"
)

const (
	idColumn               = "id"
	targetPrimaryKeySuffix = "_id"
)

func isTargetPrimaryKey(column string) bool {
	cLen := len(column)
	tLen := len(targetPrimaryKeySuffix)
	return cLen >= tLen && column[cLen-tLen:] == targetPrimaryKeySuffix
}

func isSimilarColumn(left, right database.Schema) bool {
	maybeRightTable := inflection.Plural(left.Column)
	return maybeRightTable == right.Table && right.Column == idColumn
}

// Recongnize a column thats same name of other table's primary key is a foreign key
// This base idea refers to SchemaSpy DbAnalyzer:
//   https://github.com/schemaspy/schemaspy/blob/master/src/main/java/org/schemaspy/DbAnalyzer.java
func GuessByPrimaryKeyName(schemas []database.Schema, primaryKeys []database.PrimaryKey) []constraint.Constraint {
	var constraints []constraint.Constraint

	for _, s := range schemas {
		for _, pk := range primaryKeys {
			if s.Table != pk.Table && s.Column == pk.Column && isTargetPrimaryKey(s.Column) {
				constraints = append(constraints, constraint.Constraint{s.Table, s.Column, pk.Table, pk.Column})
			}
		}
	}

	return constraints
}

func GuessBySimilarColumn(schemas []database.Schema) []constraint.Constraint {
	var constraints []constraint.Constraint

	for _, left := range schemas {
		for _, right := range schemas {
			if left.Table != right.Table && isSimilarColumn(left, right) {
				constraints = append(constraints, constraint.Constraint{left.Table, left.Column, right.Table, right.Column})
			}
		}
	}

	return constraints
}
