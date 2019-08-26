package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type authPayloadInput struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers"`
}

type authPayload struct {
	Input *authPayloadInput `json:"input"`
}

type opaResponse struct {
	Result bool `json:"result"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/authz", authz)
	e.GET("/authz", authz)
	e.Logger.Fatal(e.Start(":8080"))
}

func authz(requestContext echo.Context) error {
	authzRequestPayload := payloadFrom(requestContext.Request())
	if isAuthorized(authzRequestPayload) {
		return requestContext.NoContent(http.StatusOK)
	}
	return requestContext.String(http.StatusUnauthorized, "Unauthorized")
}

func isAuthorized(authPayload *authPayload) bool {
	authPayloadAsJSON, errm := json.Marshal(authPayload)
	if errm != nil {
		panic(errm)
	}

	opaURL := "http://localhost:8181/v1/data/ingress/allow?partial"
	authResponse, errp := http.Post(opaURL, "application/json", bytes.NewBuffer(authPayloadAsJSON))
	if errp != nil {
		panic(errp)
	}

	body, err := ioutil.ReadAll(authResponse.Body)
	if err != nil {
		panic(err.Error())
	}

	var response opaResponse
	if errj := json.Unmarshal(body, &response); errj != nil {
		panic(errj)
	}

	return response.Result
}

func payloadFrom(request *http.Request) *authPayload {
	return &authPayload{
		Input: &authPayloadInput{
			Method:  request.Header.Get("X-Original-Method"),
			Path:    request.Header.Get("X-Original-Uri"),
			Headers: addAdditionalHeaders(request.Header),
		},
	}
}

func addAdditionalHeaders(headers http.Header) http.Header {
	return headers
}
