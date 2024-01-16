package storage

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/viki-org/dnscache"

	"github.com/node-real/greenfield-bundle-service/util"
)

const (
	// OSSAccessKey defines env variable name for OSS access key
	OSSAccessId = "ALIBABA_CLOUD_ACCESS_ID"
	// OSSSecretKey defines env variable name for OSS secret key
	OSSSecretKey = "ALIBABA_CLOUD_SECRET_KEY"
	// OSSRoleARN defines env variable for OSS role arn
	OSSRoleARN = "ALIBABA_CLOUD_ROLE_ARN"
	// OSSWebIdentityTokenFile defines env variable for OSS identity token file
	OSSWebIdentityTokenFile = "ALIBABA_CLOUD_OIDC_TOKEN_FILE"
	// OSSOidcProviderArn defines env variable for OSS oidc provider arn
	OSSOidcProviderArn = "ALIBABA_CLOUD_OIDC_PROVIDER_ARN"

	// AKSKIAMType defines IAM type config which uses access key and secret key to access aws s3
	AKSKIAMType = "AKSK"
	// SAIAMType defines IAM type config which uses service account to access aws s3
	SAIAMType = "SA"
)

type OssStore struct {
	client *oss.Client
	bucket *oss.Bucket
}

type ossStorageSecretKey struct {
	accessKey string
	secretKey string
}

func getOSSSecretKeyFromEnv(accessId, secretKey string) *ossStorageSecretKey {
	key := &ossStorageSecretKey{}
	if val, ok := os.LookupEnv(accessId); ok {
		key.accessKey = val
	}
	if val, ok := os.LookupEnv(secretKey); ok {
		key.secretKey = val
	}
	return key
}

func NewOssStoreFromConfig(config *util.BundleConfig) (*OssStore, error) {
	if config.OssIAMType == AKSKIAMType {
		key := getOSSSecretKeyFromEnv(OSSAccessId, OSSSecretKey)
		return NewOssStoreAKSK(config.OssBucketUrl, key.accessKey, key.secretKey)
	} else if config.OssIAMType == SAIAMType {
		return NewOssStoreSA(config.OssBucketUrl)
	} else {
		panic("invalid oss iam type")
	}
}

func NewOssStoreSA(bucketURL string) (*OssStore, error) {
	endpoint, bucketName, _, err := parseOSS(bucketURL)
	if err != nil {
		util.Logger.Errorf("parse oss bucket error, bucketURL=%s, err=%s", bucketURL, err.Error())
		return nil, err
	}

	credProvider, err := newOIDCCredentialProvider()
	if err != nil {
		util.Logger.Errorf("failed to new oidc credential provider, %s", err)
	}
	cli, err := oss.New(endpoint, "", "", oss.SetCredentialsProvider(credProvider))
	if err != nil {
		util.Logger.Errorf("failed to use sa iam type to new oss, %s", err)
		return nil, err
	}

	bucket, err := cli.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("cannot get bucket instance %s: %s", bucketName, err)
	}

	cli.Config.Timeout = 10
	cli.Config.RetryTimes = 3
	cli.Config.HTTPTimeout.ConnectTimeout = time.Second * 30
	cli.Config.HTTPTimeout.ReadWriteTimeout = time.Second * 60
	cli.Config.HTTPTimeout.HeaderTimeout = time.Second * 60
	cli.Config.HTTPTimeout.LongTimeout = time.Second * 300
	cli.Config.IsEnableCRC = false
	cli.Config.UserAgent = "GREENFIELD-BUNDLE-SERVICE"

	return &OssStore{
		client: cli,
		bucket: bucket,
	}, nil
}

func NewOssStoreFromEnv(bucketURL string) (*OssStore, error) {
	key := getOSSSecretKeyFromEnv(OSSAccessId, OSSSecretKey)
	return NewOssStoreAKSK(bucketURL, key.accessKey, key.secretKey)
}

func checkOIDCAvailable() (bool, string, string, string) {
	oidc := true
	roleArn, exists := os.LookupEnv(OSSRoleARN)
	if !exists {
		oidc = false
		util.Logger.Error("failed to read oss role arn")
	}
	oidcProviderArn, exists := os.LookupEnv(OSSOidcProviderArn)
	if !exists {
		oidc = false
		util.Logger.Error("failed to read oss oidc provider arn")
	}
	tokenPath, exists := os.LookupEnv(OSSWebIdentityTokenFile)
	if !exists {
		oidc = false
		util.Logger.Error("failed to read oss web identity token file")
	}
	return oidc, roleArn, oidcProviderArn, tokenPath
}

type ossCredentials struct {
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
}

type defaultCredentialsProvider struct {
	cred credentials.Credential
}

func (o *ossCredentials) GetAccessKeyID() string {
	return o.AccessKeyID
}

func (o *ossCredentials) GetAccessKeySecret() string {
	return o.AccessKeySecret
}

func (o *ossCredentials) GetSecurityToken() string {
	return o.SecurityToken
}

func (d *defaultCredentialsProvider) GetCredentials() oss.Credentials {
	accessKey, _ := d.cred.GetAccessKeyId()
	secretKey, _ := d.cred.GetAccessKeySecret()
	securityToken, _ := d.cred.GetSecurityToken()
	return &ossCredentials{
		AccessKeyID:     *accessKey,
		AccessKeySecret: *secretKey,
		SecurityToken:   *securityToken,
	}
}

func newOSSCredentialsProvider(credential credentials.Credential) defaultCredentialsProvider {
	return defaultCredentialsProvider{cred: credential}
}

func newOIDCCredentialProvider() (oss.CredentialsProvider, error) {
	ok, roleArn, oidcProviderArn, tokenPath := checkOIDCAvailable()
	if !ok {
		util.Logger.Error("failed to check oss oidc")
		return nil, fmt.Errorf("no oidc env vars")
	}
	config := new(credentials.Config).
		SetType("oidc_role_arn").
		SetRoleArn(roleArn).
		SetOIDCProviderArn(oidcProviderArn).
		SetOIDCTokenFilePath(tokenPath).SetRoleSessionName("sp-oss")
	oidcCredential, err := credentials.NewCredential(config)
	if err != nil {
		util.Logger.Errorf("failed to new oidc credentials, %s", err)
	}
	provider := newOSSCredentialsProvider(oidcCredential)
	return &provider, nil
}

func NewOssStoreAKSK(bucketURL string, accessKeyId string, accessKeySecret string) (*OssStore, error) {
	endpoint, bucketName, region, err := parseOSS(bucketURL)
	if err != nil {
		util.Logger.Errorf("parse oss bucket error, bucketURL=%s, err=%s", bucketURL, err.Error())
		return nil, err
	}

	cli, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.Region(region), oss.HTTPClient(getHTTPClient(false)))
	if err != nil {
		util.Logger.Errorf("create oss client error, err=%s", err.Error())
		return nil, err
	}

	cli.Config.Timeout = 10
	cli.Config.RetryTimes = 3
	cli.Config.HTTPTimeout.ConnectTimeout = time.Second * 30
	cli.Config.HTTPTimeout.ReadWriteTimeout = time.Second * 60
	cli.Config.HTTPTimeout.HeaderTimeout = time.Second * 60
	cli.Config.HTTPTimeout.LongTimeout = time.Second * 300
	cli.Config.IsEnableCRC = false
	cli.Config.UserAgent = "GREENFIELD-BUNDLE-SERVICE"

	bucket, err := cli.Bucket(bucketName)
	if err != nil {
		return nil, fmt.Errorf("cannot get bucket instance %s: %s", bucketName, err)
	}

	return &OssStore{
		client: cli,
		bucket: bucket,
	}, nil
}

func getHTTPClient(tlsInsecureSkipVerify bool) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			// #nosec
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify},
			TLSHandshakeTimeout:   time.Second * 20,
			ResponseHeaderTimeout: time.Second * 30,
			IdleConnTimeout:       time.Second * 300,
			MaxIdleConnsPerHost:   5000,
			DialContext:           dialContext,
			DisableCompression:    true,
		},
		Timeout: time.Hour,
	}
}

func dialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	resolver := dnscache.New(time.Minute)
	rand.New(rand.NewSource(time.Now().Unix()))

	separator := strings.LastIndex(address, ":")
	if separator == -1 {
		return nil, fmt.Errorf("invalid address: %s", address)
	}
	host := address[:separator]
	port := address[separator:]
	ips, err := resolver.Fetch(host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no such host: %s", host)
	}

	var conn net.Conn
	n := len(ips)
	first := rand.Intn(n)
	dialer := &net.Dialer{Timeout: time.Second * 10}
	for i := 0; i < n; i++ {
		ip := ips[(first+i)%n]
		address = ip.String()
		if port != "" {
			address = net.JoinHostPort(address, port[1:])
		}
		conn, err = dialer.DialContext(ctx, network, address)
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
}

func parseOSS(bucketURL string) (string, string, string, error) {
	if !strings.Contains(bucketURL, "://") {
		bucketURL = fmt.Sprintf("https://%s", bucketURL)
	}
	uri, err := url.ParseRequestURI(bucketURL)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid bucket: %s, error: %v", bucketURL, err)
	}

	hostParts := strings.SplitN(uri.Host, ".", 2)
	var endpoint string
	if len(hostParts) > 1 {
		endpoint = uri.Scheme + "://" + hostParts[1]
	} else {
		return "", "", "", fmt.Errorf("cannot get oss domain name: %s", bucketURL)
	}
	regionParts := strings.SplitN(hostParts[1], ".", 2)
	if len(regionParts) != 2 {
		return "", "", "", fmt.Errorf("cannot get oss region: %s", bucketURL)
	}
	region := regionParts[0]
	bucketName := hostParts[0]

	return endpoint, bucketName, region, nil
}

func (o *OssStore) String() string {
	return fmt.Sprintf("oss://%s/", o.bucket.BucketName)
}

func (o *OssStore) GetObject(ctx context.Context, key string, off, limit int64) (resp io.ReadCloser, err error) {
	var respHeader http.Header
	if off > 0 || limit > 0 {
		var r string
		if limit > 0 {
			r = fmt.Sprintf("%d-%d", off, off+limit-1)
		} else {
			r = fmt.Sprintf("%d-", off)
		}
		resp, err = o.bucket.GetObject(key, oss.NormalizedRange(r), oss.RangeBehavior("standard"), oss.GetResponseHeader(&respHeader))
	} else {
		resp, err = o.bucket.GetObject(key, oss.GetResponseHeader(&respHeader))
		if err == nil {
			resp = verifyChecksum(resp,
				resp.(*oss.Response).Headers.Get(oss.HTTPHeaderOssMetaPrefix+ChecksumAlgo))
		}
	}
	return resp, err
}

func (o *OssStore) PutObject(ctx context.Context, key string, in io.Reader) error {
	var (
		option     []oss.Option
		respHeader http.Header
	)
	if rs, ok := in.(io.ReadSeeker); ok {
		option = append(option, oss.Meta(ChecksumAlgo, generateChecksum(rs)))
	}
	option = append(option, oss.GetResponseHeader(&respHeader))
	err := o.bucket.PutObject(key, in, option...)
	return err
}

func (o *OssStore) DeleteObject(ctx context.Context, key string) error {
	return o.bucket.DeleteObject(key)
}

func IsNoSuchKey(err error) bool {
	if err == nil {
		return false
	}
	var e oss.ServiceError
	if errors.As(err, &e) {
		return e.Code == "NoSuchKey"
	}
	return false
}
