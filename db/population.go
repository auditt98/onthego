package db

import (
	"fmt"
	"strings"
	"unicode"

	"gorm.io/gorm/clause"
)

func toSnakeCase(input string) string {
	// Remove leading and trailing whitespace
	input = strings.TrimSpace(input)

	// Replace spaces and punctuation characters with underscores
	var result strings.Builder
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || unicode.IsPunct(r) {
			result.WriteRune('_')
		}
	}

	// Convert the result to lowercase and return
	return strings.ToLower(result.String())
}

func (qe *QueryEngine) Populate(params QueryParams) *QueryEngine {
	populate := params.Populate
	hasPopulateAll := false

	for _, field := range populate {
		if field == "*" {
			hasPopulateAll = true
		} else {
			fmt.Println(toSnakeCase(field))
			// pluralize := pluralize.NewClient()
			// qe.Ref = qe.Ref.Joins("INNER JOIN " + toSnakeCase(field) + " ON " + toSnakeCase(qe.TableName) + ".id = " + toSnakeCase(field) + "." + toSnakeCase(pluralize.Singular(qe.TableName)) + "_id")
			qe.Ref = qe.Ref.Preload(field)
		}
	}
	if hasPopulateAll {
		qe.Ref = qe.Ref.Preload(clause.Associations)
	}

	return qe
}
