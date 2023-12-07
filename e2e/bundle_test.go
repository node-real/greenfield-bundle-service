package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/node-real/greenfield-bundle-service/util"
)

func TestCreateBundle(t *testing.T) {
	PrepareBundleAccounts("../cmd/bundle-service-server/db.sqlite3", 1)

	privateKey, _, err := util.GenerateRandomAccount()
	require.NoError(t, err)

	url := "http://localhost:8080/v1/createBundle"

	headers := map[string]string{
		"X-Bundle-Bucket-Name":      "example-bucket",
		"X-Bundle-Name":             "example-bundle",
		"X-Bundle-Expiry-Timestamp": fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
	}

	resp, err := SendRequest(privateKey, url, "POST", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ReadResponseBody(resp)
	require.NoError(t, err)

	t.Log("Response status:", resp.Status)
	t.Log("Response body:", body)
}

func TestFinalizeBundle(t *testing.T) {
	PrepareBundleAccounts("../cmd/bundle-service-server/db.sqlite3", 1)

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

	resp, err := SendRequest(privateKey, url, "POST", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ReadResponseBody(resp)
	require.NoError(t, err)

	t.Log("Response status:", resp.Status)
	t.Log("Response body:", body)

	// finalize bundle
	url = "http://localhost:8080/v1/finalizeBundle"

	resp, err = SendRequest(privateKey, url, "POST", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err = ReadResponseBody(resp)
	require.NoError(t, err)

	t.Log("Response status:", resp.Status)
	t.Log("Response body:", body)
}
