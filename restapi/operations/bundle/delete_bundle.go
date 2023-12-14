// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteBundleHandlerFunc turns a function with the right signature into a delete bundle handler
type DeleteBundleHandlerFunc func(DeleteBundleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteBundleHandlerFunc) Handle(params DeleteBundleParams) middleware.Responder {
	return fn(params)
}

// DeleteBundleHandler interface for that can handle valid delete bundle params
type DeleteBundleHandler interface {
	Handle(DeleteBundleParams) middleware.Responder
}

// NewDeleteBundle creates a new http.Handler for the delete bundle operation
func NewDeleteBundle(ctx *middleware.Context, handler DeleteBundleHandler) *DeleteBundle {
	return &DeleteBundle{Context: ctx, Handler: handler}
}

/*
	DeleteBundle swagger:route POST /deleteBundle Bundle deleteBundle

# Delete an bundle after object deletion on Greenfield

Delete an bundle after object deletion on Greenfield
*/
type DeleteBundle struct {
	Context *middleware.Context
	Handler DeleteBundleHandler
}

func (o *DeleteBundle) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewDeleteBundleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}