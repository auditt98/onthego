package db

type SearchParams struct {
	Filters  map[string]any `json:"$filters"`
	Page     int            `json:"$page"`
	PerPage  int            `json:"$per_page"`
	Sort     []string       `json:"$sort"`
	Populate []string       `json:"$populate"`
}
