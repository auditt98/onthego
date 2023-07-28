package db

import "strings"

func (qe *QueryEngine) Sort(params QueryParams) *QueryEngine {
	orderBy := params.OrderBy
	// order by will have the format of: "+[field]" or "-[field]"
	// example: ["+name", "-age"]
	// + means ascending, - means descending
	// if no + or - is specified, it will be ascending by default
	clauses := []string{}

	for _, field := range orderBy {
		trimmedField := strings.TrimSpace(field)
		order := "asc"
		if trimmedField[0] == '-' {
			order = "desc"
		}
		trimmedField = trimmedField[1:]
		clauses = append(clauses, trimmedField+" "+order)
	}
	sortStr := strings.Join(clauses, ", ")
	qe.Ref = qe.Ref.Order(sortStr)
	return qe
}
