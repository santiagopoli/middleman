package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/santiagopoli/middleman/internal/authorizer"
)

type ServerConfig struct {
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

	router.Post("/authz", authorizeRequest)
	router.Get("/authz", authorizeRequest)

	fmt.Println("Server Listening on :8080")
	http.ListenAndServe(":8080", router)
}
