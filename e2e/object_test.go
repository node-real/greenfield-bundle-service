package e2e

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/node-real/greenfield-bundle-service/util"

	"github.com/node-real/greenfield-bundle-service/types"
)

func uploadObject(fileName, bucketName, contentType string, file *os.File) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		return err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return err
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "http://localhost:8080/v1/uploadObject", body)
	if err != nil {
		return err
	}

	// Add headers to the request body
	headers := map[string]string{
		"X-Bundle-Bucket-Name":      bucketName,
		"X-Bundle-File-Name":        fileName,
		"X-Bundle-Content-Type":     contentType,
		"X-Bundle-Expiry-Timestamp": fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()), // Replace with the actual timestamp
	}

	for key, value := range headers {
		writer.WriteField(key, value)
	}

	// Set the content type, this is important
	req.Header.Set("Content-Type", writer.FormDataContentType())

	messageToSign := types.GetMsgToSignInBundleAuth(req)
	messageHash := types.TextHash(messageToSign)

	privateKey, _, err := util.GenerateRandomAccount()

	signature, err := SignMessage(privateKey, messageHash)
	if err != nil {
		return err
	}
	req.Header.Set(types.HTTPHeaderAuthorization, hex.EncodeToString(signature))

	// Do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// print the response body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(buf.String())

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	return nil
}

func TestUploadObject(t *testing.T) {
	file, err := os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.WriteString("test")
	if err != nil {
		t.Fatal(err)
	}

	// Upload the file
	err = uploadObject("test.txt", "test", "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}
}
