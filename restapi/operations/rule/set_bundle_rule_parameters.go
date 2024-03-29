// Code generated by go-swagger; DO NOT EDIT.

package rule

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

// NewSetBundleRuleParams creates a new SetBundleRuleParams object
//
// There are no default values defined in the spec.
func NewSetBundleRuleParams() SetBundleRuleParams {

	return SetBundleRuleParams{}
}

// SetBundleRuleParams contains all the bound params for the set bundle rule operation
// typically these are obtained from a http.Request
//
// swagger:parameters setBundleRule
type SetBundleRuleParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*User's digital signature for authorization
	  Required: true
	  In: header
	*/
	Authorization string
	/*Name of the bucket for which the rule applies
	  Required: true
	  In: header
	*/
	XBundleBucketName string
	/*Expiry timestamp of the request
	  Required: true
	  In: header
	*/
	XBundleExpiryTimestamp int64
	/*Maximum number of files in a bundle
	  Required: true
	  In: header
	*/
	XBundleMaxBundleFiles int64
	/*Maximum size of a bundle in bytes
	  Required: true
	  In: header
	*/
	XBundleMaxBundleSize int64
	/*Maximum time in seconds before a bundle must be finalized
	  Required: true
	  In: header
	*/
	XBundleMaxFinalizeTime int64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewSetBundleRuleParams() beforehand.
func (o *SetBundleRuleParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
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

	if err := o.bindXBundleMaxBundleFiles(r.Header[http.CanonicalHeaderKey("X-Bundle-Max-Bundle-Files")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleMaxBundleSize(r.Header[http.CanonicalHeaderKey("X-Bundle-Max-Bundle-Size")], true, route.Formats); err != nil {
		res = append(res, err)
	}

	if err := o.bindXBundleMaxFinalizeTime(r.Header[http.CanonicalHeaderKey("X-Bundle-Max-Finalize-Time")], true, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindAuthorization binds and validates parameter Authorization from header.
func (o *SetBundleRuleParams) bindAuthorization(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *SetBundleRuleParams) bindXBundleBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *SetBundleRuleParams) bindXBundleExpiryTimestamp(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindXBundleMaxBundleFiles binds and validates parameter XBundleMaxBundleFiles from header.
func (o *SetBundleRuleParams) bindXBundleMaxBundleFiles(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Max-Bundle-Files", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Max-Bundle-Files", "header", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("X-Bundle-Max-Bundle-Files", "header", "int64", raw)
	}
	o.XBundleMaxBundleFiles = value

	return nil
}

// bindXBundleMaxBundleSize binds and validates parameter XBundleMaxBundleSize from header.
func (o *SetBundleRuleParams) bindXBundleMaxBundleSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Max-Bundle-Size", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Max-Bundle-Size", "header", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("X-Bundle-Max-Bundle-Size", "header", "int64", raw)
	}
	o.XBundleMaxBundleSize = value

	return nil
}

// bindXBundleMaxFinalizeTime binds and validates parameter XBundleMaxFinalizeTime from header.
func (o *SetBundleRuleParams) bindXBundleMaxFinalizeTime(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("X-Bundle-Max-Finalize-Time", "header", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true

	if err := validate.RequiredString("X-Bundle-Max-Finalize-Time", "header", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("X-Bundle-Max-Finalize-Time", "header", "int64", raw)
	}
	o.XBundleMaxFinalizeTime = value

	return nil
}
