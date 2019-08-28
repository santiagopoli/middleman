package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mbndr/figlet4go"
	"github.com/santiagopoli/middleman/internal/authorizer"
	"github.com/santiagopoli/middleman/internal/http"
)

func main() {
	authorizer := authorizer.NewOPAAuthorizer("localhost:8181", "ingress/allow", true)
	authorizeRequest := http.AuthorizeRequest(authorizer)
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.POST("/authz", authorizeRequest)
	e.GET("/authz", authorizeRequest)

	PrintBanner()
	e.Logger.Fatal(e.Start(":8080"))
}

func PrintBanner() {

	ascii := figlet4go.NewAsciiRender()

	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
		figlet4go.ColorYellow,
		figlet4go.ColorCyan,
		figlet4go.ColorMagenta,
	}

	renderStr, _ := ascii.RenderOpts("Middleman", options)
	fmt.Println(renderStr)
	fmt.Println("Made with ♥ by @santiagopoli")
	fmt.Println("")
}
