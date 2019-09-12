package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/santiagopoli/middleman/internal/authorizer"
)

type ServerConfig struct {
	Port             string
	MiddlewareConfig *MiddlewareConfig
	OPAConfig        *OPAConfig
}

type MiddlewareConfig struct {
	HostHeader   string
	MethodHeader string
	PathHeader   string
}

type OPAConfig struct {
	Host                 string
	DefaultPolicy        string
	UsePartialEvaluation bool
}

func authorizeRequest(middlewareConfig *MiddlewareConfig, authorizer authorizer.Authorizer) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		authzRequestPayload := payloadFrom(request, middlewareConfig)
		if authorizer.IsAuthorized(authzRequestPayload) {
			response.WriteHeader(http.StatusOK)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
			response.Header().Add("Cache-Control", "no-cache, no-store")
			response.Write([]byte("Unauthorized"))
		}
	}
}

func payloadFrom(request *http.Request, middlewareConfig *MiddlewareConfig) *authorizer.Request {
	return &authorizer.Request{
		Host:    request.Header.Get(middlewareConfig.HostHeader),
		Method:  request.Header.Get(middlewareConfig.MethodHeader),
		Path:    request.Header.Get(middlewareConfig.PathHeader),
		Headers: addAdditionalHeaders(request.Header),
	}
}

func addAdditionalHeaders(headers http.Header) http.Header {
	return headers
}

func StartServer(serverConfig *ServerConfig) {
	authorizer := authorizer.NewOPAAuthorizer(
		serverConfig.OPAConfig.Host,
		serverConfig.OPAConfig.DefaultPolicy,
		serverConfig.OPAConfig.UsePartialEvaluation,
	)
	authorizeRequest := authorizeRequest(serverConfig.MiddlewareConfig, authorizer)
	router := chi.NewRouter()

	router.Post("/", authorizeRequest)
	router.Get("/", authorizeRequest)

	fmt.Println("Server listening on", serverConfig.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", serverConfig.Port), router)
}
