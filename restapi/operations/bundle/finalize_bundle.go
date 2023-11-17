// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// FinalizeBundleHandlerFunc turns a function with the right signature into a finalize bundle handler
type FinalizeBundleHandlerFunc func(FinalizeBundleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn FinalizeBundleHandlerFunc) Handle(params FinalizeBundleParams) middleware.Responder {
	return fn(params)
}

// FinalizeBundleHandler interface for that can handle valid finalize bundle params
type FinalizeBundleHandler interface {
	Handle(FinalizeBundleParams) middleware.Responder
}

// NewFinalizeBundle creates a new http.Handler for the finalize bundle operation
func NewFinalizeBundle(ctx *middleware.Context, handler FinalizeBundleHandler) *FinalizeBundle {
	return &FinalizeBundle{Context: ctx, Handler: handler}
}

/*
	FinalizeBundle swagger:route POST /finalizeBundle Bundle finalizeBundle

# Finalize an Existing Bundle

Completes the lifecycle of an existing bundle, requiring the bundle name for authorization.
*/
type FinalizeBundle struct {
	Context *middleware.Context
	Handler FinalizeBundleHandler
}

func (o *FinalizeBundle) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewFinalizeBundleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// FinalizeBundleBody finalize bundle body
//
// swagger:model FinalizeBundleBody
type FinalizeBundleBody struct {

	// The name of the bucket
	// Required: true
	BucketName *string `json:"bucketName"`

	// The name of the bundle to be finalized
	BundleName string `json:"bundleName,omitempty"`

	// Timestamp of the request
	// Required: true
	Timestamp *int64 `json:"timestamp"`
}

// Validate validates this finalize bundle body
func (o *FinalizeBundleBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateBucketName(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTimestamp(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *FinalizeBundleBody) validateBucketName(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"bucketName", "body", o.BucketName); err != nil {
		return err
	}

	return nil
}

func (o *FinalizeBundleBody) validateTimestamp(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"timestamp", "body", o.Timestamp); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this finalize bundle body based on context it is used
func (o *FinalizeBundleBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *FinalizeBundleBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *FinalizeBundleBody) UnmarshalBinary(b []byte) error {
	var res FinalizeBundleBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
