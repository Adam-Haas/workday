package cmd

import (
	"github.com/spf13/cobra"
	"workday/internal/worklog"
)

const defaultDateFormat = "02 Jan 06 15:04"

var wlHandler worklog.Handler

var rootCmd = &cobra.Command{
	Use: "workday",
}

func init() {
	rootCmd.AddCommand(startCommand())
	rootCmd.AddCommand(finishCommand())
	rootCmd.AddCommand(statusCommand())
}

func Execute(handler worklog.Handler) error {
	wlHandler = handler
	return rootCmd.Execute()
}