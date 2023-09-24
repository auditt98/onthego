package utils

import (
	"github.com/auditt98/onthego/types"
)

func GetIntrospection(value any, e bool) *types.IntrospectionResult {
	if e != true {
		return nil
	}
	introspection, _ := value.(*types.IntrospectionResult)
	return introspection
}
