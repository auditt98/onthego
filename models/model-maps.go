package models

import "reflect"

var ModelMap = map[string]reflect.Type{
	"user":    reflect.TypeOf(User{}),
	"article": reflect.TypeOf(Article{}),
	// Add other models as needed
}
