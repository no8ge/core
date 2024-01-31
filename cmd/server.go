package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/no8ge/core/internal/router"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

func server(port string, mod string) {
	if mod == "debug" {
		gin.SetMode(gin.DebugMode)
	}
	if mod == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	router.V1(r)
	if err := r.RunTLS(fmt.Sprintf(":%s", port), certFile, keyFile); err != nil {
		log.Fatalf("failed execute core service, %s", err.Error())
		os.Exit(1)
	}
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run core as a service",
	Long:  `Run core as a service`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		mod, _ := cmd.Flags().GetString("mod")
		server(port, mod)
	},
}
