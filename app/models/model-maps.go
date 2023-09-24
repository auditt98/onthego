package models

import "reflect"

var ModelMap = map[string]reflect.Type{
	"user":    reflect.TypeOf(User{}),
	"album":   reflect.TypeOf(Album{}),
	"photo":   reflect.TypeOf(Photo{}),
	"comment": reflect.TypeOf(Comment{}),
	// Add other models as needed
}
