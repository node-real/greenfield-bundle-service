// Code generated by go-swagger; DO NOT EDIT.

package object_upload

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// UploadObjectsHandlerFunc turns a function with the right signature into a upload objects handler
type UploadObjectsHandlerFunc func(UploadObjectsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn UploadObjectsHandlerFunc) Handle(params UploadObjectsParams) middleware.Responder {
	return fn(params)
}

// UploadObjectsHandler interface for that can handle valid upload objects params
type UploadObjectsHandler interface {
	Handle(UploadObjectsParams) middleware.Responder
}

// NewUploadObjects creates a new http.Handler for the upload objects operation
func NewUploadObjects(ctx *middleware.Context, handler UploadObjectsHandler) *UploadObjects {
	return &UploadObjects{Context: ctx, Handler: handler}
}

/*
	UploadObjects swagger:route POST /uploadObjects Object Upload uploadObjects

# Upload a single object to a bundle

Allows users to upload a single file along with a signature for validation, and a timestamp.
*/
type UploadObjects struct {
	Context *middleware.Context
	Handler UploadObjectsHandler
}

func (o *UploadObjects) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewUploadObjectsParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// UploadObjectsOKBody upload objects o k body
//
// swagger:model UploadObjectsOKBody
type UploadObjectsOKBody struct {

	// The name of the bundle where the file has been uploaded
	BundleName string `json:"bundleName,omitempty"`
}

// Validate validates this upload objects o k body
func (o *UploadObjectsOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this upload objects o k body based on context it is used
func (o *UploadObjectsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *UploadObjectsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UploadObjectsOKBody) UnmarshalBinary(b []byte) error {
	var res UploadObjectsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
