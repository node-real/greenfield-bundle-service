// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/node-real/greenfield-bundle-service/restapi/operations"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_configuration"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_file_retrieval"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_management"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/object_upload"
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
	// object_upload.UploadObjectsMaxParseMemory = 32 << 20

	if api.BundleConfigurationAddBundleRuleHandler == nil {
		api.BundleConfigurationAddBundleRuleHandler = bundle_configuration.AddBundleRuleHandlerFunc(func(params bundle_configuration.AddBundleRuleParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_configuration.AddBundleRule has not yet been implemented")
		})
	}
	if api.BundleFileRetrievalBundleFileHandler == nil {
		api.BundleFileRetrievalBundleFileHandler = bundle_file_retrieval.BundleFileHandlerFunc(func(params bundle_file_retrieval.BundleFileParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_file_retrieval.BundleFile has not yet been implemented")
		})
	}
	if api.BundleManagementManageBundleHandler == nil {
		api.BundleManagementManageBundleHandler = bundle_management.ManageBundleHandlerFunc(func(params bundle_management.ManageBundleParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_management.ManageBundle has not yet been implemented")
		})
	}
	if api.ObjectUploadUploadObjectsHandler == nil {
		api.ObjectUploadUploadObjectsHandler = object_upload.UploadObjectsHandlerFunc(func(params object_upload.UploadObjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation object_upload.UploadObjects has not yet been implemented")
		})
	}

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
