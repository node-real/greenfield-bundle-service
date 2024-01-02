// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/gin-gonic/gin"

	"github.com/bnb-chain/greenfield-go-sdk/client"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"

	"github.com/node-real/greenfield-bundle-service/auth"

	"github.com/node-real/greenfield-bundle-service/storage"

	"github.com/node-real/greenfield-bundle-service/dao"
	"github.com/node-real/greenfield-bundle-service/database"
	"github.com/node-real/greenfield-bundle-service/restapi/handlers"
	"github.com/node-real/greenfield-bundle-service/restapi/operations"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/bundle"
	"github.com/node-real/greenfield-bundle-service/restapi/operations/rule"
	"github.com/node-real/greenfield-bundle-service/service"
	"github.com/node-real/greenfield-bundle-service/util"
)

//go:generate swagger generate server --target ../../greenfield-bundle-service --name BundleService --spec ../swagger.yaml --principal interface{}

var cliOpts = struct {
	ConfigFilePath string `short:"c" long:"config-path" description:"Config path" default:"config/server/dev.json"`
}{}

func configureFlags(api *operations.BundleServiceAPI) {
	param := swag.CommandLineOptionsGroup{
		ShortDescription: "config",
		Options:          &cliOpts,
	}
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{param}
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

	api.RuleSetBundleRuleHandler = rule.SetBundleRuleHandlerFunc(handlers.HandleSetBundleRule())

	api.BundleQueryBundleHandler = bundle.QueryBundleHandlerFunc(handlers.HandleQueryBundle())

	api.BundleQueryBundlingBundleHandler = bundle.QueryBundlingBundleHandlerFunc(handlers.HandleQueryBundlingBundle())

	api.BundleViewBundleObjectHandler = bundle.ViewBundleObjectHandlerFunc(handlers.HandleViewBundleObject())

	api.BundleDownloadBundleObjectHandler = bundle.DownloadBundleObjectHandlerFunc(handlers.HandleDownloadBundleObject())

	api.BundleCreateBundleHandler = bundle.CreateBundleHandlerFunc(handlers.HandleCreateBundle())

	api.BundleDeleteBundleHandler = bundle.DeleteBundleHandlerFunc(handlers.HandleDeleteBundle())

	api.BundleFinalizeBundleHandler = bundle.FinalizeBundleHandlerFunc(handlers.HandleFinalizeBundle())

	api.BundleUploadObjectHandler = bundle.UploadObjectHandlerFunc(handlers.HandleUploadObject())

	api.BundleBundlerAccountHandler = bundle.BundlerAccountHandlerFunc(handlers.HandleGetUserBundlerAccount())

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	router := gin.Default()

	// Define the route
	router.GET("/v1/view/:bucketName/:bundleName/*objectName", func(c *gin.Context) {
		bucketName := c.Param("bucketName")
		bundleName := c.Param("bundleName")
		objectName := c.Param("objectName")
		objectName = strings.TrimPrefix(objectName, "/")

		params := bundle.NewViewBundleObjectParams()
		params.HTTPRequest = c.Request
		params.BucketName = bucketName
		params.BundleName = bundleName
		params.ObjectName = objectName

		responder := api.BundleViewBundleObjectHandler.Handle(params)

		// Write the response
		responder.WriteResponse(c.Writer, runtime.JSONProducer())
	})

	router.GET("/v1/download/:bucketName/:bundleName/*objectName", func(c *gin.Context) {
		bucketName := c.Param("bucketName")
		bundleName := c.Param("bundleName")
		objectName := c.Param("objectName")
		objectName = strings.TrimPrefix(objectName, "/")

		params := bundle.NewDownloadBundleObjectParams()
		params.HTTPRequest = c.Request
		params.BucketName = bucketName
		params.BundleName = bundleName
		params.ObjectName = objectName

		responder := api.BundleDownloadBundleObjectHandler.Handle(params)

		// Write the response
		responder.WriteResponse(c.Writer, runtime.JSONProducer())
	})

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/v1/view/") || strings.HasPrefix(req.URL.Path, "/v1/download/") {
			router.ServeHTTP(w, req)
		} else {
			setupGlobalMiddleware(api.Serve(setupMiddlewares)).ServeHTTP(w, req)
		}
	})
	return handler
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
	configFilePath := cliOpts.ConfigFilePath
	config := util.ParseServerConfigFromFile(configFilePath)

	util.InitLogger(config.LogConfig)

	db, err := database.ConnectDBWithConfig(config.DBConfig)
	if err != nil {
		panic(err)
	}

	bundleDao := dao.NewBundleDao(db)
	bundleRuleDao := dao.NewBundleRuleDao(db)
	objectDao := dao.NewObjectDao(db)
	userBundlerAccountDao := dao.NewUserBundlerAccountDao(db)
	bundlerAccountDao := dao.NewBundlerAccountDao(db)

	gnfdClient, err := client.New(config.GnfdConfig.ChainId, config.GnfdConfig.RpcUrl, client.Option{})
	if err != nil {
		panic(fmt.Errorf("unable to new greenfield client, %v", err))
	}

	// set a random default account for server gnfd client
	privkey, _, err := util.GenerateRandomAccount()
	if err != nil {
		panic(err)
	}
	serverAccount, err := types.NewAccountFromPrivateKey("server-account", hex.EncodeToString(privkey))
	if err != nil {
		panic(err)
	}
	util.Logger.Infof("set greenfield client default server account: %s", serverAccount.GetAddress().String())
	gnfdClient.SetDefaultAccount(serverAccount)

	fileManager := storage.NewFileManager(config, objectDao, gnfdClient)
	authManager := auth.NewAuthManager(gnfdClient)

	// init services
	service.GnfdClient = gnfdClient

	service.BundleSvc = service.NewBundleService(gnfdClient, authManager, bundleDao, bundleRuleDao, userBundlerAccountDao)
	service.BundleRuleSvc = service.NewBundleRuleService(bundleRuleDao)
	service.ObjectSvc = service.NewObjectService(config, fileManager, bundleDao, objectDao, userBundlerAccountDao)
	service.UserBundlerAccountSvc = service.NewUserBundlerAccountService(userBundlerAccountDao, bundlerAccountDao)
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
