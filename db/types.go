package db

import (
	"encoding/json"

	"gorm.io/gorm"
)

type QueryEngine struct {
	// FindOne       func(params QueryParams) (interface{}, error)
	// FindMany      func(params QueryParams) ([]interface{}, error)
	// FindWithCount func(params QueryParams) ([]interface{}, int64, error)
	// Create        func(modelName string, data interface{}) (interface{}, error)
	// Update        func(modelName string, id uint, data interface{}) (interface{}, error)
	// Delete        func(modelName string, id uint) (interface{}, error)
	Ref *gorm.DB
}

type QueryParams struct {
	Picks    []string
	Where    WhereParams
	Offset   int
	OrderBy  []string
	Limit    int
	Populate []string //user.articles, user.articles.comments
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
	Not  []interface{}
	Attr map[string]AttributeParams
}

// var a = QueryParams{
// 	Where: WhereParams{
// 		Or: []WhereParams{
// 			{
// 				Attr: map[string]AttributeParams{
// 					"email": {
// 						Eq: "test",
// 					},
// 				},
// 			},
// 			{
// 				Attr: map[string]AttributeParams{
// 					"email": {
// 						Eq: "test",
// 					},
// 				},
// 			},
// 		},
// 		Attr: map[string]AttributeParams{
// 			"email": {
// 				Eq: "test",
// 			},
// 		},
// 	},
// }
