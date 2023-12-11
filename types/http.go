package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	HTTPHeaderContentType = "Content-Type"
	HTTPHeaderUnsignedMsg = "X-Bundle-Unsigned-Msg"

	HTTPHeaderFileSHA256      = "X-Bundle-File-Sha256"
	HTTPHeaderBucketName      = "X-Bundle-Bucket-Name"
	HTTPHeaderAttributes      = "X-Bundle-Attributes"
	HTTPHeaderMaxBundleSize   = "X-Bundle-Max-Bundle-Size"
	HTTPHeaderMaxFileSize     = "X-Bundle-Max-File-Size"
	HTTPHeaderMaxFinalizeTime = "X-Bundle-Max-Finalize-Time"

	// HTTPHeaderExpiryTimestamp defines the expiry timestamp, which is the ISO 8601 datetime string (e.g. 2021-09-30T16:25:24Z), and the maximum Timestamp since the request sent must be less than MaxExpiryAgeInSec (seven days).
	HTTPHeaderExpiryTimestamp = "X-Bundle-Expiry-Timestamp"
	HTTPHeaderAuthorization   = "Authorization"

	// MaxExpiryAgeInSec defines the maximum expiry age in seconds
	MaxExpiryAgeInSec = 3600 * 24 * 7 // 7 days
)

var supportedHeaders = []string{
	HTTPHeaderFileSHA256,
	HTTPHeaderContentType,
	HTTPHeaderUnsignedMsg,
	HTTPHeaderBucketName,
	HTTPHeaderAttributes,
	HTTPHeaderMaxBundleSize,
	HTTPHeaderMaxFileSize,
	HTTPHeaderMaxFinalizeTime,
	HTTPHeaderExpiryTimestamp,
}

func initSupportHeaders() map[string]struct{} {
	supportMap := make(map[string]struct{})
	for _, header := range supportedHeaders {
		emptyStruct := new(struct{})
		supportMap[header] = *emptyStruct
	}
	return supportMap
}

// EncodePath encode the strings from UTF-8 byte representations to HTML hex escape sequences
func EncodePath(pathName string) string {
	reservedNames := regexp.MustCompile("^[a-zA-Z0-9-_.~/]+$")
	// no need to encode
	if reservedNames.MatchString(pathName) {
		return pathName
	}
	var encodedPathName strings.Builder
	for _, s := range pathName {
		if 'A' <= s && s <= 'Z' || 'a' <= s && s <= 'z' || '0' <= s && s <= '9' { // §2.3 Unreserved characters (mark)
			encodedPathName.WriteRune(s)
			continue
		}
		switch s {
		case '-', '_', '.', '~', '/':
			encodedPathName.WriteRune(s)
			continue
		default:
			length := utf8.RuneLen(s)
			if length < 0 {
				// if utf8 cannot convert return the same string as is
				return pathName
			}
			u := make([]byte, length)
			utf8.EncodeRune(u, s)
			for _, r := range u {
				hexStr := hex.EncodeToString([]byte{r})
				encodedPathName.WriteString("%" + strings.ToUpper(hexStr))
			}
		}
	}
	return encodedPathName.String()
}

// GetHostInfo returns host header from the request
func GetHostInfo(req *http.Request) string {
	host := req.Header.Get("host")
	if host != "" {
		return host
	}
	if req.Host != "" {
		return req.Host
	}
	return req.URL.Host
}

// getSignedHeaders return the sorted header array
func getSortedHeaders(req *http.Request, supportMap map[string]struct{}) []string {
	var signHeaders []string
	for k := range req.Header {
		if _, ok := supportMap[k]; ok {
			signHeaders = append(signHeaders, strings.ToLower(k))
		}
	}
	sort.Strings(signHeaders)
	return signHeaders
}

// getCanonicalHeaders generate a list of request headers with their values
func getCanonicalHeaders(req *http.Request, supportHeaders map[string]struct{}) string {
	var content bytes.Buffer
	var containHostHeader bool
	sortHeaders := getSortedHeaders(req, supportHeaders)
	headerMap := make(map[string][]string)
	for key, data := range req.Header {
		headerMap[strings.ToLower(key)] = data
	}

	for _, header := range sortHeaders {
		content.WriteString(strings.ToLower(header))
		content.WriteByte(':')

		if header != "host" {
			for i, v := range headerMap[header] {
				if i > 0 {
					content.WriteByte(',')
				}
				trimVal := strings.Join(strings.Fields(v), " ")
				content.WriteString(trimVal)
			}
			content.WriteByte('\n')
		} else {
			containHostHeader = true
			content.WriteString(GetHostInfo(req))
			content.WriteByte('\n')
		}
	}

	if !containHostHeader {
		content.WriteString(GetHostInfo(req))
		content.WriteByte('\n')
	}
	return content.String()
}

// getSignedHeaders return the alphabetically sorted, semicolon-separated list of lowercase request header names.
func getSignedHeaders(req *http.Request, supportHeaders map[string]struct{}) string {
	return strings.Join(getSortedHeaders(req, supportHeaders), ";")
}

// GetCanonicalRequest generate the canonicalRequest base on aws s3 sign without payload hash.
func GetCanonicalRequest(req *http.Request) string {
	supportHeaders := initSupportHeaders()
	req.URL.RawQuery = strings.ReplaceAll(req.URL.Query().Encode(), "+", "%20")
	canonicalRequest := strings.Join([]string{
		req.Method,
		EncodePath(req.URL.Path),
		req.URL.RawQuery,
		getCanonicalHeaders(req, supportHeaders),
		getSignedHeaders(req, supportHeaders),
	}, "\n")
	return canonicalRequest
}

func GetMsgToSignInBundleAuth(req *http.Request) []byte {
	return crypto.Keccak256([]byte(GetCanonicalRequest(req)))
}

// TextHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//
//	keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func TextHash(data []byte) []byte {
	hash, _ := TextAndHash(data)
	return hash
}

// TextAndHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calculated as
//
//	keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

// VerifySignature verifies the signature of the given message hash and returns the signer's address.
func VerifySignature(req *http.Request) (common.Address, error) {
	messageToSign := GetMsgToSignInBundleAuth(req)
	messageHash := TextHash(messageToSign)

	requestSignature := req.Header.Get(HTTPHeaderAuthorization)

	sigBytes, err := hex.DecodeString(requestSignature)
	if err != nil {
		return common.Address{}, err
	}

	address, err := util.RecoverAddress(common.BytesToHash(messageHash), sigBytes)
	if err != nil {
		return common.Address{}, err
	}
	return address, nil
}

// ValidateExpiryTimestamp validates the expiry timestamp
func ValidateExpiryTimestamp(req *http.Request) error {
	expiryTimestamp := req.Header.Get(HTTPHeaderExpiryTimestamp)
	if expiryTimestamp == "" {
		return fmt.Errorf("expiry timestamp is empty")
	}

	// parse expiry timestamp from int
	expiryTime, err := strconv.ParseInt(expiryTimestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("expiry timestamp is invalid")
	}

	if expiryTime < time.Now().Unix() {
		return fmt.Errorf("expiry timestamp is expired")
	}
	if expiryTime-time.Now().Unix() > MaxExpiryAgeInSec {
		return fmt.Errorf("expiry timestamp is too far in the future")
	}
	return nil
}
