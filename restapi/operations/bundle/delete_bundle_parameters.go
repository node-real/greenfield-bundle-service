// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewDeleteBundleParams creates a new DeleteBundleParams object
//
// There are no default values defined in the spec.
func NewDeleteBundleParams() DeleteBundleParams {

	return DeleteBundleParams{}
}

// DeleteBundleParams contains all the bound params for the delete bundle operation
// typically these are obtained from a http.Request
//
// swagger:parameters deleteBundle
type DeleteBundleParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*User's digital signature for authorization
	  Required: true
	  In: header
	*/
	Authorization string
	/*The name of the bucket
	  Required: true
	  In: header
	*/
	XBundleBucketName string
	/*Expiry timestamp of the request
	  Required: true
	  In: header
	*/
	XBundleExpiryTimestamp int64
	/*The name of the bundle to be finalized
	  Required: true
	  In: header
	*/
	XBundleName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewDeleteBundleParams() beforehand.
func (o *DeleteBundleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if err := o.bindAuthorization(r.Header[http.CanonicalHeaderKey("Authorization")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleBucketName(r.Header[http.CanonicalHeaderKey("X-Bundle-Bucket-Name")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleExpiryTimestamp(r.Header[http.CanonicalHeaderKey("X-Bundle-Expiry-Timestamp")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleName(r.Header[http.CanonicalHeaderKey("X-Bundle-Name")], true, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *DeleteBundleParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *DeleteBundleParams) bindXBundleBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindXBundleExpiryTimestamp binds and validates parameter XBundleExpiryTimestamp from header.
func (o *DeleteBundleParams) bindXBundleExpiryTimestamp(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindXBundleName binds and validates parameter XBundleName from header.
func (o *DeleteBundleParams) bindXBundleName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Name", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Name", "header", raw); err != nil {
		return err
	}
	o.XBundleName = raw

	return nil
}
