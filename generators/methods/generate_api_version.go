package generators

import (
	"fmt"
	"os"
	"path/filepath"
)

func GenerateAPIVersionFolder(rootpath string, path string) {
	// Concatenate the root path and the specified path
	absPath := filepath.Join(rootpath, path)

	// Check if the folder already exists
	_, err := os.Stat(absPath)
	if err == nil {
		fmt.Printf("Folder '%s' already exists\n", absPath)
		return
	}

	// Create the folder
	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Failed to create folder '%s': %v\n", absPath, err)
		return
	}

	fmt.Printf("Created folder '%s'\n", absPath)
}

func GenerateAPIVersionRouter() {
	fmt.Println("Generate API version router")
}
