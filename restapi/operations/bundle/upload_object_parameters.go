// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// UploadObjectMaxParseMemory sets the maximum size in bytes for
// the multipart form parser for this operation.
//
// The default value is 32 MB.
// The multipart parser stores up to this + 10MB.
var UploadObjectMaxParseMemory int64 = 32 << 20

// NewUploadObjectParams creates a new UploadObjectParams object
//
// There are no default values defined in the spec.
func NewUploadObjectParams() UploadObjectParams {

	return UploadObjectParams{}
}

// UploadObjectParams contains all the bound params for the upload object operation
// typically these are obtained from a http.Request
//
// swagger:parameters uploadObject
type UploadObjectParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*User's digital signature for authentication
	  Required: true
	  In: header
	*/
	Authorization string
	/*The name of the bucket
	  Required: true
	  In: header
	*/
	XBundleBucketName string
	/*Content type of the file
	  Required: true
	  In: header
	*/
	XBundleContentType string
	/*Expiry timestamp of the request
	  Required: true
	  In: header
	*/
	XBundleExpiryTimestamp int64
	/*The name of the file to be uploaded
	  Required: true
	  In: header
	*/
	XBundleFileName string
	/*The file to be uploaded
	  Required: true
	  In: formData
	*/
	File io.ReadCloser
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewUploadObjectParams() beforehand.
func (o *UploadObjectParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := r.ParseMultipartForm(UploadObjectMaxParseMemory); err != nil {
		if err != http.ErrNotMultipart {
			return errors.New(400, "%v", err)
		} else if err := r.ParseForm(); err != nil {
			return errors.New(400, "%v", err)
		}
	}

	if err := o.bindAuthorization(r.Header[http.CanonicalHeaderKey("Authorization")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleBucketName(r.Header[http.CanonicalHeaderKey("X-Bundle-Bucket-Name")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleContentType(r.Header[http.CanonicalHeaderKey("X-Bundle-Content-Type")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleExpiryTimestamp(r.Header[http.CanonicalHeaderKey("X-Bundle-Expiry-Timestamp")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleFileName(r.Header[http.CanonicalHeaderKey("X-Bundle-File-Name")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		res = append(res, errors.New(400, "reading file %q failed: %v", "file", err))
	} else if err := o.bindFile(file, fileHeader); err != nil {
		// Required: true
		res = append(res, err)
	} else {
		o.File = &runtime.File{Data: file, Header: fileHeader}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *UploadObjectParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("Authorization", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("Authorization", "header", raw); err != nil {
		return err
	}
	o.Authorization = raw

	return nil
}

// bindXBundleBucketName binds and validates parameter XBundleBucketName from header.
func (o *UploadObjectParams) bindXBundleBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Bucket-Name", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Bucket-Name", "header", raw); err != nil {
		return err
	}
	o.XBundleBucketName = raw

	return nil
}

// bindXBundleContentType binds and validates parameter XBundleContentType from header.
func (o *UploadObjectParams) bindXBundleContentType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Content-Type", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Content-Type", "header", raw); err != nil {
		return err
	}
	o.XBundleContentType = raw

	return nil
}

// bindXBundleExpiryTimestamp binds and validates parameter XBundleExpiryTimestamp from header.
func (o *UploadObjectParams) bindXBundleExpiryTimestamp(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Expiry-Timestamp", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Expiry-Timestamp", "header", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("X-Bundle-Expiry-Timestamp", "header", "int64", raw)
	}
	o.XBundleExpiryTimestamp = value

	return nil
}

// bindXBundleFileName binds and validates parameter XBundleFileName from header.
func (o *UploadObjectParams) bindXBundleFileName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-File-Name", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-File-Name", "header", raw); err != nil {
		return err
	}
	o.XBundleFileName = raw

	return nil
}

// bindFile binds file parameter File.
//
// The only supported validations on files are MinLength and MaxLength
func (o *UploadObjectParams) bindFile(file multipart.File, header *multipart.FileHeader) error {
	return nil
}
