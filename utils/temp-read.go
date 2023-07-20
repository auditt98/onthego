package utils

import (
	"fmt"
	"html/template"
	"io/ioutil"
)

func ReadTemplate(path string) (*template.Template, error) {
	templateContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read template file '%s': %v\n", path, err)
		return nil, err
	}

	// Create a new template
	tmpl := template.New("handler")
	tmpl, err = tmpl.Parse(string(templateContent))
	if err != nil {
		fmt.Printf("Failed to parse template: %v\n", err)
		return nil, err
	}
	return tmpl, nil
}
