package middlewares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func QueryAnalyzerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryParams := c.Request.URL.Query()

		queryObject := make(map[string]interface{})
		var query_params_list []string
		query_params_list = make([]string, 0)

		for key, _ := range queryParams {
			if strings.HasPrefix(key, "filters") {
				parts := strings.Split(key, "[")
				param := parts[1:]
				for i, _ := range param {
					param[i] = strings.TrimSuffix(param[i], "]")
				}
				query_params_list = append(query_params_list, param...)
				//append the value of the query
				query_params_list = append(query_params_list, queryParams[key][0])
			}
		}
		c.JSON(200, query_params_list)
		queryObject["filters"] = createNestedObject(query_params_list)
		c.Set("queryObject", queryObject)
		c.Next()
	}
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func createNestedObject(arr []string) map[string]interface{} {
	root := make(map[string]interface{})
	// current := root
	// prevCurrent := root
	validOperations := []string{
		"$eq", "$neq", "$gt", "$gte", "$lt", "$lte", "$in",
		"$contains", "$startswith",
		"$endswith",
	}
	queries := [][]string{}

	for i := 0; i < len(arr)-1; i++ {
		item := arr[i]
		// nextItem := arr[i+1]
		if item == "$or" {
			//look for the next valid operation, add from $or to that operation + 1 to another array
			for j := i + 1; j < len(arr); j++ {
				if contains(validOperations, arr[j]) {
					addArr := arr[i : j+1]
					addArr = append(addArr, arr[j+1])
					// addArr = remove(arr, 1)
					queries = append(queries, addArr)
					//remove the added array from the original array
					i = j + 1
					break
				}
			}
		} else {
			for j := i + 1; j < len(arr); j++ {
				if contains(validOperations, arr[j]) {
					addArr := arr[i : j+1]
					addArr = append(addArr, arr[j+1])
					// addArr = remove(arr, 1)
					queries = append(queries, addArr)
					//remove the added array from the original array
					i = j
					break
				}
			}
		}
		// if item == "$or" {

		// 	if _, ok := root["$or"]; !ok {
		// 		root["$or"] = make(map[string]interface{})
		// 	}
		// 	current = root["$or"].(map[string]interface{})
		// 	current[nextItem] = make(map[string]interface{})
		// 	prevCurrent = current
		// 	current = current[nextItem].(map[string]interface{})
		// 	i++
		// } else {
		// 	if contains(validOperations, item) {
		// 		var key string
		// 		for k := range prevCurrent {
		// 			key = k
		// 			break
		// 		}
		// 		switch item {
		// 		case "$eq":
		// 			prevCurrent[key] = arr[i+1]
		// 		case "$neq":
		// 			prevCurrent[key] = "!=" + arr[i+1]
		// 		case "$gt":
		// 			prevCurrent[key] = ">" + arr[i+1]
		// 		case "$gte":
		// 			prevCurrent[key] = ">=" + arr[i+1]
		// 		case "$lt":
		// 			prevCurrent[key] = "<" + arr[i+1]
		// 		case "$lte":
		// 			prevCurrent[key] = "<=" + arr[i+1]
		// 		case "$in":
		// 			prevCurrent[key] = arr[i+1]
		// 		case "$contains":
		// 			prevCurrent[key] = "~%" + arr[i+1] + "%"
		// 		case "$startswith":
		// 			prevCurrent[key] = "~" + arr[i+1] + "%"
		// 		case "$endswith":
		// 			prevCurrent[key] = "~" + "%" + arr[i+1]
		// 		}
		// 		prevCurrent = current
		// 		current = root
		// 		i++
		// 		continue
		// 	}
		// 	current[item] = make(map[string]interface{})
		// 	prevCurrent = current
		// 	current = current[item].(map[string]interface{})
		// }
	}
	fmt.Println("queries: ", queries)
	return root
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
