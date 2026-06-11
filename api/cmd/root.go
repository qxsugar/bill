package cmd

import (
	"github.com/qxsugar/bill/api/cmd/http"
	"github.com/qxsugar/bill/api/cmd/worker"
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

// workerCmd 启动后台 worker，消费 asynq 异步任务（如房间超时自动结算）。
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "start background worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.Setup()
		w, cleanup, err := worker.InitializeWorker()
		if err != nil {
			return err
		}
		defer cleanup()
		w.Start()
		return nil
	},
}

func Execute() {
	rootCmd.AddCommand(workerCmd)
	cobra.CheckErr(rootCmd.Execute())
}
