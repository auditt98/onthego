package db

import (
	"encoding/json"

	"gorm.io/gorm"
)

type QueryEngine struct {
	Ref       *gorm.DB
	TableName string
}

type QueryParams struct {
	Picks    []string
	Where    WhereParams
	Offset   int
	OrderBy  []string
	Limit    int
	Populate []string //*, user.articles, user.articles.comments
}

type AttributeParams struct {
	Eq          interface{}
	Eqi         interface{}
	Ne          interface{}
	In          []interface{}
	Nin         []interface{}
	Lt          json.Number
	Lte         json.Number
	Gt          json.Number
	Gte         json.Number
	Between     []json.Number
	Contains    string
	NContains   string
	Containsi   string
	NContainsi  string
	StartsWith  string
	EndsWith    string
	NStartsWith string
	Nil         bool
	NNil        bool
}

type WhereParams struct {
	And  []WhereParams
	Or   []WhereParams
	Attr map[string]AttributeParams
}
