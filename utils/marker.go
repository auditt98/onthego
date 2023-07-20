package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReplaceMarkerInSection(marker, replaceString, filePath string) error {
	// Get the absolute path of the root directory
	rootDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("failed to get root directory path: %v", err)
	}

	// Concatenate the root directory path with the provided file path
	absoluteFilePath := filepath.Join(rootDir, filePath)

	// Read the contents of the file
	content, err := ioutil.ReadFile(absoluteFilePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Find the start and end markers
	startMarker := fmt.Sprintf("//generator: %s", marker)
	endMarker := fmt.Sprintf("//endgenerator: %s", marker)

	// Find the indices of the start and end markers
	startIndex := strings.Index(string(content), startMarker)
	endIndex := strings.Index(string(content), endMarker)

	fmt.Println("------Start, End: ", string(content), endIndex)

	// Ensure both start and end markers exist in the file
	if startIndex == -1 || endIndex == -1 {
		return fmt.Errorf("failed to find markers in the file")
	}

	// Calculate the replacement position
	startPos := startIndex + len(startMarker)
	endPos := endIndex

	// Generate the new content by replacing the section between markers with the replaceString
	newContent := string(content[:startPos]) + "\n" + replaceString + "\n" + string(content[endPos:])

	// Write the modified content back to the file
	err = ioutil.WriteFile(absoluteFilePath, []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	cmd := exec.Command("gofmt", "-w", absoluteFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run gofmt: %v", err)
	}
	return nil
}
