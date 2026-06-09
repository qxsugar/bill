package cmd

import (
	"github.com/qxsugar/bill/api/cmd/http"
	"github.com/qxsugar/bill/api/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "bill-api",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Setup()
		app, cleanup, err := http.InitializeApplication()
		if err != nil {
			return err
		}
		defer cleanup()
		app.Start()
		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
