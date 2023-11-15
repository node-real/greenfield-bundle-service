// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_configuration"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_file_retrieval"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle_management"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/object_upload"
)

// NewBundleServiceAPI creates a new BundleService instance
func NewBundleServiceAPI(spec *loads.Document) *BundleServiceAPI {
	return &BundleServiceAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		useSwaggerUI:        false,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer:          runtime.JSONConsumer(),
		MultipartformConsumer: runtime.DiscardConsumer,

		BinProducer: runtime.ByteStreamProducer(),

		BundleConfigurationAddBundleRuleHandler: bundle_configuration.AddBundleRuleHandlerFunc(func(params bundle_configuration.AddBundleRuleParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_configuration.AddBundleRule has not yet been implemented")
		}),
		BundleFileRetrievalBundleFileHandler: bundle_file_retrieval.BundleFileHandlerFunc(func(params bundle_file_retrieval.BundleFileParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_file_retrieval.BundleFile has not yet been implemented")
		}),
		BundleManagementManageBundleHandler: bundle_management.ManageBundleHandlerFunc(func(params bundle_management.ManageBundleParams) middleware.Responder {
			return middleware.NotImplemented("operation bundle_management.ManageBundle has not yet been implemented")
		}),
		ObjectUploadUploadObjectsHandler: object_upload.UploadObjectsHandlerFunc(func(params object_upload.UploadObjectsParams) middleware.Responder {
			return middleware.NotImplemented("operation object_upload.UploadObjects has not yet been implemented")
		}),
	}
}

/*BundleServiceAPI API for handling file bundling and querying objects in the Bundle Service. */
type BundleServiceAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler
	useSwaggerUI    bool

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator

	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator

	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer
	// MultipartformConsumer registers a consumer for the following mime types:
	//   - multipart/form-data
	MultipartformConsumer runtime.Consumer

	// BinProducer registers a producer for the following mime types:
	//   - application/octet-stream
	BinProducer runtime.Producer

	// BundleConfigurationAddBundleRuleHandler sets the operation handler for the add bundle rule operation
	BundleConfigurationAddBundleRuleHandler bundle_configuration.AddBundleRuleHandler
	// BundleFileRetrievalBundleFileHandler sets the operation handler for the bundle file operation
	BundleFileRetrievalBundleFileHandler bundle_file_retrieval.BundleFileHandler
	// BundleManagementManageBundleHandler sets the operation handler for the manage bundle operation
	BundleManagementManageBundleHandler bundle_management.ManageBundleHandler
	// ObjectUploadUploadObjectsHandler sets the operation handler for the upload objects operation
	ObjectUploadUploadObjectsHandler object_upload.UploadObjectsHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// UseRedoc for documentation at /docs
func (o *BundleServiceAPI) UseRedoc() {
	o.useSwaggerUI = false
}

// UseSwaggerUI for documentation at /docs
func (o *BundleServiceAPI) UseSwaggerUI() {
	o.useSwaggerUI = true
}

// SetDefaultProduces sets the default produces media type
func (o *BundleServiceAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *BundleServiceAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *BundleServiceAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *BundleServiceAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *BundleServiceAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *BundleServiceAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *BundleServiceAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the BundleServiceAPI
func (o *BundleServiceAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}
	if o.MultipartformConsumer == nil {
		unregistered = append(unregistered, "MultipartformConsumer")
	}

	if o.BinProducer == nil {
		unregistered = append(unregistered, "BinProducer")
	}

	if o.BundleConfigurationAddBundleRuleHandler == nil {
		unregistered = append(unregistered, "bundle_configuration.AddBundleRuleHandler")
	}
	if o.BundleFileRetrievalBundleFileHandler == nil {
		unregistered = append(unregistered, "bundle_file_retrieval.BundleFileHandler")
	}
	if o.BundleManagementManageBundleHandler == nil {
		unregistered = append(unregistered, "bundle_management.ManageBundleHandler")
	}
	if o.ObjectUploadUploadObjectsHandler == nil {
		unregistered = append(unregistered, "object_upload.UploadObjectsHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *BundleServiceAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *BundleServiceAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *BundleServiceAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *BundleServiceAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		case "multipart/form-data":
			result["multipart/form-data"] = o.MultipartformConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *BundleServiceAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/octet-stream":
			result["application/octet-stream"] = o.BinProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *BundleServiceAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the bundle service API
func (o *BundleServiceAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *BundleServiceAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/addBundleRule"] = bundle_configuration.NewAddBundleRule(o.context, o.BundleConfigurationAddBundleRuleHandler)
	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/bundleFile/{bundleName}/{objectName}"] = bundle_file_retrieval.NewBundleFile(o.context, o.BundleFileRetrievalBundleFileHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/manageBundle"] = bundle_management.NewManageBundle(o.context, o.BundleManagementManageBundleHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/uploadObjects"] = object_upload.NewUploadObjects(o.context, o.ObjectUploadUploadObjectsHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *BundleServiceAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	if o.useSwaggerUI {
		return o.context.APIHandlerSwaggerUI(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *BundleServiceAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *BundleServiceAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *BundleServiceAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *BundleServiceAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[um][path] = builder(h)
	}
}
