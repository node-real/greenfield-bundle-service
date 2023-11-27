// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/node-real/greenfield-bundle-service/restapi/handlers"

	"github.com/node-real/greenfield-bundle-service/restapi/operations"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"
)

//go:generate swagger generate server --target ../../greenfield-bundle-service --name BundleService --spec ../swagger.yaml --principal interface{}

func configureFlags(api *operations.BundleServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.BundleServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()
	api.MultipartformConsumer = runtime.DiscardConsumer

	api.BinProducer = runtime.ByteStreamProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// bundle.UploadObjectMaxParseMemory = 32 << 20

	api.RuleSetBundleRuleHandler = rule.SetBundleRuleHandlerFunc(handlers.HandleAddBundleRule())

	api.BundleBundleObjectHandler = bundle.BundleObjectHandlerFunc(handlers.HandleBundleObject())

	api.BundleCreateBundleHandler = bundle.CreateBundleHandlerFunc(handlers.HandleCreateBundle())

	api.BundleFinalizeBundleHandler = bundle.FinalizeBundleHandlerFunc(handlers.HandleFinalizeBundle())

	api.BundleUploadObjectHandler = bundle.UploadObjectHandlerFunc(handlers.HandleUploadObject())

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
	// initialize service.BundleSvc
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
