package guess

import (
	"testing"

	"github.com/syucream/hakagi/src/database"
)

func TestGuessByPrimaryKey(t *testing.T) {
	cases := []struct{
		left database.Column
		rightTable string
		right database.Column
		expected bool
	}{
		// matching cases
		{
			database.Column{
				Name: "example_id",
				Type: "int",
			},
			"examples",
			database.Column{
				Name: "example_id",
				Type: "int",
			},
			true,
		},

		// mis-matching cases
		{
			database.Column{
				Name: "foo_id",
				Type: "int",
			},
			"examples",
			database.Column{
				Name: "bar_id",
				Type: "int",
			},
			false,
		},
		{
			database.Column{
				Name: "foo_id",
				Type: "text",
			},
			"examples",
			database.Column{
				Name: "bar_id",
				Type: "text",
			},
			false,
		},
		{
			database.Column{
				Name: "foo_id",
				Type: "blob",
			},
			"examples",
			database.Column{
				Name: "bar_id",
				Type: "blob",
			},
			false,
		},
	}

	guesser := GuessByPrimaryKey()
	for _, c := range cases {
		if actual := guesser(c.left, c.rightTable, c.right); actual != c.expected {
			t.Errorf("Actual: %v, Expected: %v\n", actual, c.expected)
		}
	}
}

func TestGuessByTableAndColumn(t *testing.T) {
	cases := []struct{
		left database.Column
		rightTable string
		right database.Column
		expected bool
	}{
		// matching cases
		{
			database.Column{
				Name: "example_id",
				Type: "int",
			},
			"examples",
			database.Column{
				Name: "id",
				Type: "int",
			},
			true,
		},

		// mis-matching cases
		{
			database.Column{
				Name: "id",
				Type: "int",
			},
			"examples",
			database.Column{
				Name: "id",
				Type: "int",
			},
			false,
		},
		{
			database.Column{
				Name: "example_id",
				Type: "text",
			},
			"examples",
			database.Column{
				Name: "id",
				Type: "text",
			},
			false,
		},
		{
			database.Column{
				Name: "example_id",
				Type: "blob",
			},
			"examples",
			database.Column{
				Name: "id",
				Type: "blob",
			},
			false,
		},
	}

	guesser := GuessByTableAndColumn()
	for _, c := range cases {
		if actual := guesser(c.left, c.rightTable, c.right); actual != c.expected {
			t.Errorf("Actual: %v, Expected: %v\n", actual, c.expected)
		}
	}
}
