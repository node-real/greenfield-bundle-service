package e2e

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

func RandomString(length int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 生成随机字符
	var sb strings.Builder
	for i := 0; i < length; i++ {
		index := rand.Intn(len(alphabet))
		sb.WriteByte(alphabet[index])
	}

	return sb.String()
}

func TestCreateBundle(t *testing.T) {
	privateKey, _, err := util.GenerateRandomAccount()
	require.NoError(t, err)

	url := "http://localhost:8080/v1/createBundle"

	headers := map[string]string{
		"X-Bundle-Bucket-Name":      "example-bucket",
		"X-Bundle-Name":             "example-bundle",
		"X-Bundle-Expiry-Timestamp": fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
	}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	messageToSign := types.GetMsgToSignInBundleAuth(req)
	messageHash := types.TextHash(messageToSign)

	signature, err := SignMessage(privateKey, messageHash)
	require.NoError(t, err)

	req.Header.Set(types.HTTPHeaderAuthorization, hex.EncodeToString(signature))

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))
}

func TestFinalizeBundle(t *testing.T) {
	privateKey, addr, err := util.GenerateRandomAccount()
	println(addr.String())
	require.NoError(t, err)

	bundleName := RandomString(10)
	bucketName := RandomString(10)

	headers := map[string]string{
		"X-Bundle-Bucket-Name":      bucketName,
		"X-Bundle-Name":             bundleName,
		"X-Bundle-Expiry-Timestamp": fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
	}

	url := "http://localhost:8080/v1/createBundle"

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	messageToSign := types.GetMsgToSignInBundleAuth(req)
	messageHash := types.TextHash(messageToSign)

	signature, err := SignMessage(privateKey, messageHash)
	require.NoError(t, err)

	req.Header.Set(types.HTTPHeaderAuthorization, hex.EncodeToString(signature))

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", hex.EncodeToString(signature))

	// Create a new HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))

	// finalize bundle
	url = "http://localhost:8080/v1/finalizeBundle"

	// Create a new request
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	messageToSign = types.GetMsgToSignInBundleAuth(req)
	messageHash = types.TextHash(messageToSign)

	signature, err = SignMessage(privateKey, messageHash)
	require.NoError(t, err)

	req.Header.Set(types.HTTPHeaderAuthorization, hex.EncodeToString(signature))

	// Create a new HTTP client and send the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))
}
