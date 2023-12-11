package e2e

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func uploadObject(privateKey []byte, fileName, bucketName, contentType string, file *os.File) error {
	// Create a new SHA256 hash
	hash := sha256.New()

	// Write the file's content to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	// Get the hash sum in bytes
	hashInBytes := hash.Sum(nil)[:]

	// Hex encode the hash sum
	hashInHex := hex.EncodeToString(hashInBytes)

	// Reset the file read pointer to the beginning
	file.Seek(0, 0)

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

	err = writer.Close()
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Content-Type":              writer.FormDataContentType(),
		"X-Bundle-Bucket-Name":      bucketName,
		"X-Bundle-File-Name":        fileName,
		"X-Bundle-Content-Type":     contentType,
		"X-Bundle-Expiry-Timestamp": fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
		"X-Bundle-File-Sha256":      hashInHex,
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	return nil
}

func GenerateRandomHTMLPage() string {
	return "<html><body><h1>Random number: " + strconv.Itoa(10) + "</h1></body></html>"
}

func TestUploadObject(t *testing.T) {
	PrepareBundleAccounts("../cmd/bundle-service-server/db.sqlite3", 1)

	privateKey, _, err := GetAccount()
	if err != nil {
		t.Fatal(err)
	}

	// Create a temporary file
	file, err := os.CreateTemp("", "test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name()) // Clean up

	_, err = file.WriteString(GenerateRandomHTMLPage())
	if err != nil {
		t.Fatal(err)
	}
	file.Close() // Close the file after writing to it

	// Upload the file
	file, err = os.Open(file.Name()) // Open the file for reading
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close() // Ensure the file is closed after it is no longer needed

	// Upload the file
	err = uploadObject(privateKey, "test.html", "bundle-test", "text/html", file)
	if err != nil {
		t.Fatal(err)
	}
}
