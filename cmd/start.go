package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"time"
	"workday/internal/worklog"
)

func startCommand() *cobra.Command {
	return &cobra.Command{
		Use: "start",
		Short: "Records the workday starting",
		RunE: func (cmd *cobra.Command, args []string) error {
			startTime, err := wlHandler.StartNewWorkday()
			if err == worklog.ErrPreviousWorkdayIncomplete {
				return errors.New ("the previous workday is incomplete, please run `workday finish --time hh:mm to specify a finish time for the previous workday`")
			}
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Worday started at %s", startTime.Format(time.RFC1123)))
			return nil
		},
	}
}