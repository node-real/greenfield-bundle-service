// Code generated by go-swagger; DO NOT EDIT.

package bundle

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// QueryBundlingBundleHandlerFunc turns a function with the right signature into a query bundling bundle handler
type QueryBundlingBundleHandlerFunc func(QueryBundlingBundleParams) middleware.Responder

// Handle executing the request and returning a response
func (fn QueryBundlingBundleHandlerFunc) Handle(params QueryBundlingBundleParams) middleware.Responder {
	return fn(params)
}

// QueryBundlingBundleHandler interface for that can handle valid query bundling bundle params
type QueryBundlingBundleHandler interface {
	Handle(QueryBundlingBundleParams) middleware.Responder
}

// NewQueryBundlingBundle creates a new http.Handler for the query bundling bundle operation
func NewQueryBundlingBundle(ctx *middleware.Context, handler QueryBundlingBundleHandler) *QueryBundlingBundle {
	return &QueryBundlingBundle{Context: ctx, Handler: handler}
}

/*
	QueryBundlingBundle swagger:route GET /queryBundlingBundle/{bucketName} Bundle queryBundlingBundle

# Query bundling bundle information of a bucket

Queries the bundling bundle information of a given bucket.
*/
type QueryBundlingBundle struct {
	Context *middleware.Context
	Handler QueryBundlingBundleHandler
}

func (o *QueryBundlingBundle) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewQueryBundlingBundleParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
