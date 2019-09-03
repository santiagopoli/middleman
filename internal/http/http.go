package http

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/santiagopoli/middleman/internal/authorizer"
)

func authorizeRequest(authorizer authorizer.Authorizer) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		authzRequestPayload := payloadFrom(request)
		if authorizer.IsAuthorized(authzRequestPayload) {
			response.WriteHeader(http.StatusOK)
		} else {
			response.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func payloadFrom(request *http.Request) *authorizer.Request {
	return &authorizer.Request{
		Method:  request.Header.Get("X-Original-Method"),
		Path:    request.Header.Get("X-Original-Uri"),
		Headers: addAdditionalHeaders(request.Header),
	}
}

func addAdditionalHeaders(headers http.Header) http.Header {
	return headers
}

func StartServer() {
	authorizer := authorizer.NewOPAAuthorizer("localhost:8181", "ingress/allow", false)
	authorizeRequest := authorizeRequest(authorizer)
	router := chi.NewRouter()

	router.Post("/authz", authorizeRequest)
	router.Get("/authz", authorizeRequest)

	fmt.Println("Server Listening on :8080")
	http.ListenAndServe(":8080", router)
}
