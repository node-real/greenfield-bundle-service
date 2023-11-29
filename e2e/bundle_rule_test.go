package e2e

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"

	"github.com/node-real/greenfield-bundle-service/restapi/handlers"
	"github.com/node-real/greenfield-bundle-service/util"
)

// SignMessage signs a message with a given private key
func SignMessage(privateKeyBytes []byte, message []byte) ([]byte, error) {
	// Convert bytes to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Sign the message
	signature, err := crypto.Sign(crypto.Keccak256(message), privateKey)
	if err != nil {
		return nil, err
	}
	return signature, err
}

func TestSetOrCreateBundleRule(t *testing.T) {
	privateKey, _, err := util.GenerateRandomAccount()
	require.NoError(t, err)

	// sign request
	signMessage := handlers.BundleRuleSignMessage{
		Method:          handlers.SetBundleRuleMethod,
		BucketName:      "test",
		MaxFiles:        10,
		MaxSize:         1000000,
		MaxFinalizeTime: 10000,
		Timestamp:       time.Now().Unix(),
	}

	signBytes, err := signMessage.SignBytes()
	require.NoError(t, err)

	signature, err := SignMessage(privateKey, signBytes)
	require.NoError(t, err)

	res, err := util.VerifySignature(crypto.Keccak256(signBytes), signature)
	require.NoError(t, err)

	reqBody := rule.SetBundleRuleBody{
		BucketName:      &signMessage.BucketName,
		MaxBundleSize:   &signMessage.MaxFiles,
		MaxBundleFiles:  &signMessage.MaxSize,
		MaxFinalizeTime: &signMessage.MaxFinalizeTime,
		Timestamp:       &signMessage.Timestamp,
	}

	url := "http://localhost:8080/v1/setBundleRule"
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

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))

	println(res)
}
