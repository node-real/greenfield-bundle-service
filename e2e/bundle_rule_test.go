package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSetOrCreateBundleRule(t *testing.T) {
	PrepareBundleAccounts("../cmd/bundle-service-server/db.sqlite3", 1)

	privateKey, _, err := GetAccount()
	require.NoError(t, err)

	url := "http://localhost:8080/v1/setBundleRule"

	headers := map[string]string{
		"Content-Type":               "application/json",
		"X-Bundle-Bucket-Name":       "bundle-test",
		"X-Bundle-Max-Bundle-Size":   "1048576",
		"X-Bundle-Max-Bundle-Files":  "100",
		"X-Bundle-Max-Finalize-Time": "3600",
		"X-Bundle-Expiry-Timestamp":  fmt.Sprintf("%d", time.Now().Add(1*time.Hour).Unix()),
	}

	resp, err := SendRequest(privateKey, url, "POST", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := ReadResponseBody(resp)
	require.NoError(t, err)

	t.Log("Response status:", resp.Status)
	t.Log("Response body:", body)

	resp, err = SendRequest(privateKey, url, "POST", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err = ReadResponseBody(resp)
	require.NoError(t, err)

	t.Log("Response status:", resp.Status)
	t.Log("Response body:", body)
}
