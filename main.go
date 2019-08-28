package main

import (
	"fmt"
	"github.com/mbndr/figlet4go"
	"github.com/santiagopoli/middleman/internal/http"
)

func main() {
	PrintBanner()
	http.StartServer()
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
