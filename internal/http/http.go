package http

import (
	"net/http"

	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo"
	"github.com/santiagopoli/middleman/internal/authorizer"
)

func authorizeRequest(authorizer authorizer.Authorizer) func(echo.Context) error {
	return func(requestContext echo.Context) error {
		authzRequestPayload := payloadFrom(requestContext.Request())
		if authorizer.IsAuthorized(authzRequestPayload) {
			return requestContext.NoContent(http.StatusOK)
		}
		return requestContext.String(http.StatusUnauthorized, "Unauthorized")
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
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.POST("/authz", authorizeRequest)
	e.GET("/authz", authorizeRequest)

	e.Logger.Fatal(e.Start(":8080"))
}
