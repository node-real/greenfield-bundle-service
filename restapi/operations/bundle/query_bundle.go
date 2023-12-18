// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// QueryBundleHandlerFunc turns a function with the right signature into a query bundle handler
type QueryBundleHandlerFunc func(QueryBundleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn QueryBundleHandlerFunc) Handle(params QueryBundleParams) middleware.Responder {
	return fn(params)
}

// QueryBundleHandler interface for that can handle valid query bundle params
type QueryBundleHandler interface {
	Handle(QueryBundleParams) middleware.Responder
}

// NewQueryBundle creates a new http.Handler for the query bundle operation
func NewQueryBundle(ctx *middleware.Context, handler QueryBundleHandler) *QueryBundle {
	return &QueryBundle{Context: ctx, Handler: handler}
}

/*
	QueryBundle swagger:route GET /queryBundle/{bucketName}/{bundleName} Bundle queryBundle

# Query bundle information

Queries a specific object from a given bundle and returns its related information.
*/
type QueryBundle struct {
	Context *middleware.Context
	Handler QueryBundleHandler
}

func (o *QueryBundle) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewQueryBundleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
