package utils

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
)

func FileUpload(files []*multipart.FileHeader, prefix string) {

	//switch UPLOAD_DRIVER s3, local

	switch os.Getenv("UPLOAD_DRIVER") {
	case "s3":
		//upload to s3
		HandleS3FileUpload(files, prefix)
	case "local":
		HandleLocalFileUpload(files, prefix)
	}

}

func HandleS3FileUpload(files []*multipart.FileHeader, prefix string) {

}

func HandleLocalFileUpload(files []*multipart.FileHeader, prefix string) error {
	// Get the file upload path from the environment variable
	uploadPath := os.Getenv("FILE_UPLOAD_PATH")
	if uploadPath == "" {
		return fmt.Errorf("FILE_UPLOAD_PATH environment variable is not set")
	}

	// Create the directory if it doesn't exist
	err := os.MkdirAll(uploadPath+"/"+prefix, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}

	for _, file := range files {
		fileContent, err := file.Open()
		if err != nil {
			return fmt.Errorf("Error opening file: %v", err)
		}
		defer fileContent.Close()

		byteContainer, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return fmt.Errorf("Error reading file content: %v", err)
		}

		escapedFilename := url.QueryEscape(file.Filename)
		filePath := uploadPath + "/" + prefix + "/" + escapedFilename

		err = ioutil.WriteFile(filePath, byteContainer, 0644)
		if err != nil {
			return fmt.Errorf("Error writing file: %v", err)
		}
	}

	return nil
}
