package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssStore struct {
	client *oss.Client
	bucket *oss.Bucket
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
