package main

import (
	"fmt"

	"github.com/santiagopoli/middleman/internal/http"
)

func main() {
	PrintBanner()
	http.StartServer()
}

func PrintBanner() {
	fmt.Println("Middleman!")
	fmt.Println("Made with ♥ by @santiagopoli")
	fmt.Println("")
}
