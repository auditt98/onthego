package generators

import (
	"fmt"
	"strings"

	"github.com/auditt98/onthego/utils"
)

func GenerateActionString(group, action, version, templatePath string) (string, error) {
	tmpl, err := utils.ReadTemplate(templatePath)
	// Prepare the data for template execution
	data := struct {
		Group   string
		Action  string
		Version string
	}{
		Group:   strings.Title(group),
		Action:  strings.Title(action),
		Version: strings.Title(version),
	}
	var handlerContent strings.Builder
	err = tmpl.Execute(&handlerContent, data)
	if err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		return "", err
	}
	return handlerContent.String(), nil
}

func GenerateSingleAction(group, action, handlerPath, templatePath string) (string, error) {
	return "", nil
}

func GenerateMultipleAction(actions map[string]utils.RouteConfig, group, version, handlerPath, templatePath string) (string, error) {

	actionString := ""
	for _, route := range actions {
		action, _ := GenerateActionString(group, route.Handler, version, templatePath)
		actionString += action + "\n"
	}
	utils.ReplaceMarkerInSection("Actions", actionString, handlerPath)
	return actionString, nil
}

func GenerateActionRouter() {
	fmt.Println()
}
