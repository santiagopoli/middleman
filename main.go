package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/santiagopoli/middleman/internal/authorizer"
	"github.com/santiagopoli/middleman/internal/http"
)

func main() {
	authorizer := authorizer.NewOPAAuthorizer("http://localhost:8181/v1/data/ingress/allow?partial")
	authorizeRequest := http.AuthorizeRequest(authorizer)
	e := echo.New()
	e.Use(middleware.Logger())
	e.POST("/authz", authorizeRequest)
	e.GET("/authz", authorizeRequest)
	e.Logger.Fatal(e.Start(":8080"))
}
