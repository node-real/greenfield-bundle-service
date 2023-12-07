package e2e

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/node-real/greenfield-bundle-service/util"
)

func uploadObject(privateKey []byte, fileName, bucketName, contentType string, file *os.File) error {
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

	url := "http://localhost:8080/v1/uploadObject"
	resp, err := SendRequest(privateKey, url, "POST", headers, body.Bytes())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyStr, err := ReadResponseBody(resp)
	if err != nil {
		return err
	}

	fmt.Println(bodyStr)

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	return nil
}

func TestUploadObject(t *testing.T) {
	PrepareBundleAccounts("../cmd/bundle-service-server/db.sqlite3", 1)

	privateKey, _, err := util.GenerateRandomAccount()
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = file.WriteString("test")
	if err != nil {
		t.Fatal(err)
	}

	// Upload the file
	err = uploadObject(privateKey, "test.txt", "test", "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}
}
