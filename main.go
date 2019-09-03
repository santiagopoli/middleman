package main

import (
	"github.com/santiagopoli/middleman/internal/http"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "middleman",
		Short: "Authz Middleware powered by Open Policy Agent",
		Run:   startServer(),
	}
	rootCmd.Flags().String("middleware.HostHeader", "X-Original-Host", "Header to use as Host")
	rootCmd.Flags().String("middleware.MethodHeader", "X-Original-Method", "Header to use as Method")
	rootCmd.Flags().String("middleware.PathHeader", "X-Original-Uri", "Header to use as Path")
	rootCmd.Flags().String("opa.Host", "localhost:8181", "Location of the OPA Server")
	rootCmd.Flags().String("opa.DefaultPolicy", "ingress/allow", "Default Policy to Use")
	rootCmd.Flags().Bool("opa.UsePartialEvaluation", false, "Use partial evaluation for Policy checks")

	rootCmd.Execute()
}

func startServer() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		http.StartServer(buildConfigFromFlags(cmd.Flags()))
	}
}

func buildConfigFromFlags(flags *pflag.FlagSet) *http.ServerConfig {
	middlewareHostHeader, _ := flags.GetString("middleware.HostHeader")
	middlewareMethodHeader, _ := flags.GetString("middleware.MethodHeader")
	middlewarePathHeader, _ := flags.GetString("middleware.PathHeader")

	opaHost, _ := flags.GetString("opa.Host")
	opaDefaultPolicy, _ := flags.GetString("opa.DefaultPolicy")
	opaUsePartialEvaluation, _ := flags.GetBool("opa.UsePartialEvaluation")

	return &http.ServerConfig{
		MiddlewareConfig: &http.MiddlewareConfig{
			HostHeader:   middlewareHostHeader,
			MethodHeader: middlewareMethodHeader,
			PathHeader:   middlewarePathHeader,
		},
		OPAConfig: &http.OPAConfig{
			Host:                 opaHost,
			DefaultPolicy:        opaDefaultPolicy,
			UsePartialEvaluation: opaUsePartialEvaluation,
		},
	}
}
