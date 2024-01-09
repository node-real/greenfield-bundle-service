package e2e

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/bnb-chain/greenfield-bundle-sdk/bundle"
	bundleTypes "github.com/bnb-chain/greenfield-bundle-sdk/types"
)

func visit(root string, b *bundle.Bundle) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			relativePath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			contentType := mime.TypeByExtension(ext)

			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}

			hash := sha256.Sum256(content)

			options := &bundleTypes.AppendObjectOptions{
				ContentType: contentType,
				HashAlgo:    bundleTypes.HashAlgo_SHA256, // Set the hash algorithm to SHA256
				Hash:        hash[:],                     // Set the hash
			}

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				return err
			}

			_, err = b.AppendObject(relativePath, file, options)
			if err != nil {
				return err
			}

			println(relativePath, hex.EncodeToString(hash[:]))
		}
		return nil
	}
}

func bundleDirectory(dir string) (io.ReadSeekCloser, int64, error) {
	b, err := bundle.NewBundle()
	if err != nil {
		return nil, 0, err
	}

	err = filepath.Walk(dir, visit(dir, b))
	if err != nil {
		return nil, 0, err
	}

	return b.FinalizeBundle()
}

func saveBundleToFile(bundle io.ReadSeekCloser, filePath string) error {
	// Create a new file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the contents of the bundle to the file
	_, err = io.Copy(file, bundle)
	if err != nil {
		return err
	}

	return nil
}

func TestPack(t *testing.T) {
	bundleObject, _, err := bundleDirectory(".")
	if err != nil {
		t.Fatal(err)
	}
	err = saveBundleToFile(bundleObject, "./cmd.bundle")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadBundle(t *testing.T) {
	// Pack a directory into a bundle and save it to a file
	bundleObject, _, err := bundleDirectory("../Vinayak-09.github.io/")
	if err != nil {
		t.Fatal(err)
	}
	bundleFilePath := "./web.bundle"
	err = saveBundleToFile(bundleObject, bundleFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// Open the bundle file and read its content
	bundleFile, err := os.Open(bundleFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer bundleFile.Close()
	bundleContent, err := ioutil.ReadAll(bundleFile)
	if err != nil {
		t.Fatal(err)
	}

	// Calculate the SHA256 hash of the bundle file content
	hash := sha256.Sum256(bundleContent)
	hashInHex := hex.EncodeToString(hash[:])

	// Create a new multipart form and add the bundle file to it
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	filePart, err := writer.CreateFormFile("file", "bundle")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(filePart, bytes.NewReader(bundleContent))
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Set the necessary headers for the request
	headers := map[string]string{
		"Content-Type":              writer.FormDataContentType(),
		"X-Bundle-Bucket-Name":      "bundle-test",
		"X-Bundle-Name":             "mysite",
		"X-Bundle-File-Sha256":      hashInHex,
		"X-Bundle-Expiry-Timestamp": strconv.FormatInt(time.Now().Add(1*time.Hour).Unix(), 10),
	}

	privateKey, add, err := GetAccount()
	println(add.String())
	if err != nil {
		t.Fatal(err)
	}

	// Send a POST request to the uploadBundle endpoint
	url := "http://localhost:8080/v1/uploadBundle"
	resp, err := SendRequest(privateKey, url, "POST", headers, body.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("received non-OK response status: %s", resp.Status)
	}
}
