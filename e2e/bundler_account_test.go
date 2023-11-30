package e2e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/node-real/greenfield-bundle-service/util"
)

func getBundlerAccount(userAddress string) (string, error) {
	// Construct the URL with the userAddress path parameter
	url := fmt.Sprintf("http://localhost:8080/v1/bundlerAccount/%s", userAddress)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	// Set the header to application/json, as required
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request error: %s", string(body))
	}

	return string(body), nil
}

func TestGetUserBundlerAccount(t *testing.T) {
	_, address, err := util.GenerateRandomAccount()
	if err != nil {
		t.Fatal(err)
	}

	bundlerAccount, err := getBundlerAccount(address.String())
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("bundler account: %s", bundlerAccount)
}
