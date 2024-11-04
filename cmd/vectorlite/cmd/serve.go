/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"VectorLite/internal/api"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a VectorLite instance",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		api.Serve(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().Int("port", 9123, "Port to run the server on")
}
