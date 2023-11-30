package e2e

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/node-real/greenfield-bundle-service/restapi/handlers"
)

func uploadObject(signature, fileName, bucketName, bundleName, contentType string, timestamp int64, file *os.File) error {
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

	// Add other form fields
	_ = writer.WriteField("fileName", fileName)
	_ = writer.WriteField("bucketName", bucketName)
	_ = writer.WriteField("bundleName", bundleName)
	_ = writer.WriteField("timestamp", fmt.Sprintf("%d", timestamp))
	_ = writer.WriteField("contentType", contentType)

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

	// Set the content type, this is important
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Set the X-Signature header
	req.Header.Set("X-Signature", signature)

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

func signMessage(message *handlers.ObjectSignMessage) (signature []byte, address string, err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	data, err := message.SignBytes()
	if err != nil {
		return nil, "", err
	}

	hash := crypto.Keccak256Hash(data)

	signature, err = crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return nil, "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, "", fmt.Errorf("error casting public key to ECDSA")
	}

	// Derive the address from the public key
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return signature, address, nil
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

	signMsg := &handlers.ObjectSignMessage{
		Method:      handlers.UploadObjectMethod,
		BucketName:  "test",
		BundleName:  "test",
		FileName:    "test.txt",
		ContentType: "text/plain",
		Timestamp:   time.Now().Unix(),
	}
	// Sign the message
	signature, _, err := signMessage(signMsg)
	if err != nil {
		t.Fatal(err)
	}

	// Upload the file
	err = uploadObject(hex.EncodeToString(signature), "test.txt", "test", "test", "text/plain", signMsg.Timestamp, file)
	if err != nil {
		t.Fatal(err)
	}
}
