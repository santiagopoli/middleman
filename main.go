package main

import (
	"fmt"

	"github.com/santiagopoli/middleman/internal/http"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "middleman",
		Short: "Authz Middleware powered by Open Policy Agent",
		Run:   startServer(),
	}
	rootCmd.Flags().String("opa.Host", "localhost:8181", "Location of the OPA Server")
	rootCmd.Flags().String("opa.DefaultPolicy", "ingress/allow", "Default Policy to Use")
	rootCmd.Flags().Bool("opa.UsePartialEvaluation", false, "Use partial evaluation for Policy checks")
	rootCmd.Execute()
}

func startServer() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		opaHost, _ := cmd.Flags().GetString("opa.Host")
		opaDefaultPolicy, _ := cmd.Flags().GetString("opa.DefaultPolicy")
		opaUsePartialEvaluation, _ := cmd.Flags().GetBool("opa.UsePartialEvaluation")
		PrintBanner()
		http.StartServer(opaHost, opaDefaultPolicy, opaUsePartialEvaluation)
	}
}

func PrintBanner() {
	fmt.Println("Middleman!")
	fmt.Println("Made with ♥ by @santiagopoli")
	fmt.Println("")
}
