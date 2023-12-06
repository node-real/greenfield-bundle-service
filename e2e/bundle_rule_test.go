package e2e

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/node-real/greenfield-bundle-service/types"
)

// SignMessage signs a message with a given private key
func SignMessage(privateKeyBytes []byte, message []byte) ([]byte, error) {
	// Convert bytes to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	println("hex message", hex.EncodeToString(message))

	// Sign the message
	signature, err := crypto.Sign(message, privateKey)
	if err != nil {
		return nil, err
	}
	println("hex signature", hex.EncodeToString(signature))
	return signature, err
}

func GetAccount() ([]byte, common.Address, error) {
	privateKeyBytes, _ := hex.DecodeString("4feafe85242413ba7914121ecc43406de4e5199d343660190768d68f87fe8611")
	// Convert the bytes to *ecdsa.PrivateKey
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	return privateKeyBytes, crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

func TestSetOrCreateBundleRule(t *testing.T) {
	privateKey, acc, err := GetAccount()
	require.NoError(t, err)
	println("address", acc.String())

	url := "http://localhost:8080/v1/setBundleRule"

	// Define the headers based on the Swagger specification
	headers := map[string]string{
		"X-Bundle-Bucket-Name":       "t3ply5",
		"X-Bundle-Max-Bundle-Size":   "1048576", // 1 MB in bytes
		"X-Bundle-Max-Bundle-Files":  "100",
		"X-Bundle-Max-Finalize-Time": "3600", // 1 hour in seconds
		"X-Bundle-Expiry-Timestamp":  fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		panic(err)
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")

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
}
