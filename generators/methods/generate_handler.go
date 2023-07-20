package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/auditt98/onthego/utils"
)

func GenerateHandlerFile(rootpath, filename, version, templatePath string) {
	// Concatenate the root path and the specified filename
	absPath := filepath.Join(rootpath, filename+".go")

	// Check if the file already exists
	_, err := os.Stat(absPath)
	if err == nil {
		fmt.Printf("File '%s' already exists\n", absPath)
		return
	}

	tmpl, err := utils.ReadTemplate(templatePath)

	// Prepare the data for template execution
	data := struct {
		Name    string
		Actions string
		Version string
	}{
		Name:    strings.Title(filename),
		Actions: "//generator: Actions\n//endgenerator: Actions",
		Version: strings.Title(version),
	}

	// Generate the handler content by executing the template with the data
	var handlerContent strings.Builder
	err = tmpl.Execute(&handlerContent, data)
	if err != nil {
		fmt.Printf("Failed to execute template: %v\n", err)
		return
	}

	// Create the file
	file, err := os.Create(absPath)
	if err != nil {
		fmt.Printf("Failed to create file '%s': %v\n", absPath, err)
		return
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(handlerContent.String())
	if err != nil {
		fmt.Printf("Failed to write content to file '%s': %v\n", absPath, err)
		return
	}

	fmt.Printf("Created file '%s'\n", absPath)
}

func GenerateHandlerRouter() {
	fmt.Println()
}
