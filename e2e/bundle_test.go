package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/node-real/greenfield-bundle-service/restapi/handlers"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
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

	// sign request
	signMessage := handlers.BundleSignMessage{
		Method:     handlers.CreateBundleMethod,
		BucketName: "test",
		BundleName: "test",
		Timestamp:  time.Now().Unix(),
	}

	signBytes, err := signMessage.SignBytes()
	require.NoError(t, err)

	signature, err := SignMessage(privateKey, signBytes)
	require.NoError(t, err)

	reqBody := bundle.CreateBundleBody{
		BucketName: &signMessage.BucketName,
		BundleName: &signMessage.BundleName,
		Timestamp:  &signMessage.Timestamp,
	}
	url := "http://localhost:8080/v1/createBundle"
	jsonData, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

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
}

func TestFinalizeBundle(t *testing.T) {
	privateKey, _, err := util.GenerateRandomAccount()
	require.NoError(t, err)

	bundleName := RandomString(10)
	bucketName := RandomString(10)
	// sign request
	signMessage := handlers.BundleSignMessage{
		Method:     handlers.CreateBundleMethod,
		BucketName: bucketName,
		BundleName: bundleName,
		Timestamp:  time.Now().Unix(),
	}

	signBytes, err := signMessage.SignBytes()
	require.NoError(t, err)

	signature, err := SignMessage(privateKey, signBytes)
	require.NoError(t, err)

	reqBody := bundle.CreateBundleBody{
		BucketName: &signMessage.BucketName,
		BundleName: &signMessage.BundleName,
		Timestamp:  &signMessage.Timestamp,
	}
	url := "http://localhost:8080/v1/createBundle"
	jsonData, err := json.Marshal(reqBody)
	require.NoError(t, err)

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

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
	signMessage = handlers.BundleSignMessage{
		Method:     handlers.FinalizeBundleMethod,
		BucketName: bucketName,
		BundleName: bundleName,
		Timestamp:  time.Now().Unix(),
	}

	signBytes, err = signMessage.SignBytes()
	require.NoError(t, err)

	signature, err = SignMessage(privateKey, signBytes)
	require.NoError(t, err)

	finalizeReqBody := bundle.FinalizeBundleBody{
		BucketName: &signMessage.BucketName,
		BundleName: &signMessage.BundleName,
		Timestamp:  &signMessage.Timestamp,
	}
	url = "http://localhost:8080/v1/finalizeBundle"
	jsonData, err = json.Marshal(finalizeReqBody)
	require.NoError(t, err)

	// Create a new request
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Signature", hex.EncodeToString(signature))

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
