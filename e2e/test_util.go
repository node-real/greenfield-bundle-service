package e2e

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/types"
	"github.com/node-real/greenfield-bundle-service/util"
)

func GetAccount() ([]byte, common.Address, error) {
	privateKeyBytes, _ := hex.DecodeString("your private key")
	// Convert the bytes to *ecdsa.PrivateKey
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	return privateKeyBytes, crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

// SignMessage signs a message with a given private key
func SignMessage(privateKeyBytes []byte, message []byte) ([]byte, error) {
	// Convert bytes to ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Sign the message
	signature, err := crypto.Sign(message, privateKey)
	if err != nil {
		return nil, err
	}
	return signature, err
}

// PrepareBundleAccounts prepares the bundle accounts for testing
func PrepareBundleAccounts(dbPath string, n int) {
	db, err := database.ConnectDBWithConfig(&util.DBConfig{
		DBDialect: "sqlite3",
		DBPath:    dbPath,
	})
	if err != nil {
		util.Logger.Errorf("connect to db error, err=%s", err.Error())
		return
	}

	var count int64
	db.Model(&database.BundlerAccount{}).Count(&count)

	for i := count; i < int64(n); i++ {
		_, account, err := util.GenerateRandomAccount()
		if err != nil {
			util.Logger.Errorf("generate random account error, err=%s", err.Error())
			continue
		}

		bundlerAccount := database.BundlerAccount{
			AccountAddress: account.String(),
			Status:         database.BundleAccountStatusLiving,
		}

		result := db.Create(&bundlerAccount)
		if result.Error != nil {
			util.Logger.Errorf("create bundler account error, err=%s", result.Error.Error())
		}
	}
}

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

func SendRequest(privateKey []byte, url, method string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	messageToSign := types.GetMsgToSignInBundleAuth(req)
	messageHash := types.TextHash(messageToSign)

	signature, err := SignMessage(privateKey, messageHash)
	if err != nil {
		return nil, err
	}

	req.Header.Set(types.HTTPHeaderAuthorization, hex.EncodeToString(signature))

	client := &http.Client{}
	return client.Do(req)
}

func ReadResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
