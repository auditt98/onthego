package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

func FileUpload(files []*multipart.FileHeader, prefix string) ([]string, error) {

	//switch UPLOAD_DRIVER s3, local

	switch os.Getenv("UPLOAD_DRIVER") {
	case "s3":
		//upload to s3
		return HandleS3FileUpload(files, prefix)
	case "local":
		return HandleLocalFileUpload(files, prefix)
	}
	return nil, fmt.Errorf("UPLOAD_DRIVER environment variable is not set")
}

func HandleS3FileUpload(files []*multipart.FileHeader, prefix string) ([]string, error) {
	return nil, nil
}

func HandleLocalFileUpload(files []*multipart.FileHeader, prefix string) ([]string, error) {
	// Get the file upload path from the environment variable
	uploadPath := os.Getenv("FILE_UPLOAD_PATH")
	if uploadPath == "" {
		return nil, fmt.Errorf("FILE_UPLOAD_PATH environment variable is not set")
	}
	paths := []string{}
	// Create the directory if it doesn't exist
	err := os.MkdirAll(uploadPath+"/"+prefix, 0755)
	if err != nil {
		return nil, fmt.Errorf("Error creating directory: %v", err)
	}

	for _, file := range files {
		fileContent, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("Error opening file: %v", err)
		}
		defer fileContent.Close()

		byteContainer, err := ioutil.ReadAll(fileContent)
		if err != nil {
			return nil, fmt.Errorf("Error reading file content: %v", err)
		}

		escapedFilename := url.QueryEscape(file.Filename)
		filePath := uploadPath + "/" + prefix + "/" + escapedFilename
		basePath := "/" + prefix + "/" + escapedFilename
		//check if file exists
		if _, err := os.Stat(filePath); err == nil {
			// File already exists, generate a new unique file name
			fileExt := path.Ext(escapedFilename) // Get the file extension
			baseNameWithoutExt := strings.TrimSuffix(escapedFilename, fileExt)
			counter := 1
			for {
				// Generate a new file name with the counter
				newBaseName := fmt.Sprintf("%s (%d)%s", baseNameWithoutExt, counter, fileExt)
				newBaseName = url.QueryEscape(newBaseName)
				newPath := uploadPath + "/" + prefix + "/" + newBaseName
				// Check if the new file name exists
				_, err = os.Stat(newPath)
				if os.IsNotExist(err) {
					// New file name is unique, use it
					filePath = newPath
					basePath = "/" + prefix + "/" + newBaseName
					break
				}

				counter++
			}
		}
		err = ioutil.WriteFile(filePath, byteContainer, 0755)
		if err != nil {
			return nil, fmt.Errorf("Error writing file: %v", err)
		}
		paths = append(paths, basePath)
	}

	return paths, nil
}

func GeneratePresignedUrl(resourcePath, secretKey string, expiration time.Duration) string {
	// Calculate the expiration time
	expirationTime := time.Now().UTC().Add(expiration)

	// Prepare the query parameters
	queryParams := url.Values{}
	queryParams.Set("expires", fmt.Sprintf("%d", expirationTime.Unix()))

	signature := CalculateSignature(resourcePath, secretKey, expirationTime.Unix())
	// Add the signature to the query parameters
	queryParams.Set("signature", signature)

	// Construct the presigned URL
	presignedURL := fmt.Sprintf("%s?%s", resourcePath, queryParams.Encode())

	return presignedURL
}

func CalculateSignature(resourcePath, secretKey string, expiration int64) string {
	// Create the string to sign
	stringToSign := fmt.Sprintf("%s\n%d", resourcePath, expiration)
	// Calculate the signature
	hmacSha256 := hmac.New(sha256.New, []byte(secretKey))
	hmacSha256.Write([]byte(stringToSign))
	signature := base64.URLEncoding.EncodeToString(hmacSha256.Sum(nil))
	return signature
}

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()
		result[fieldName] = fieldValue
	}

	return result
}
