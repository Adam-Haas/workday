package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func statusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Outputs the current workday state and any accrued overtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := wlHandler.GetStatus()
			if err != nil {
				return err
			}

			tableData := make([][]string, len(*s.UnusedOvertimeEntries))
			for i, entry := range *s.UnusedOvertimeEntries {
				tableData[i] = []string{strconv.FormatInt(entry.ID, 10), entry.StartTime.Format(defaultDateFormat), entry.FinishTime.Format(defaultDateFormat), entry.OvertimeReason}
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Start Time", "Finish Time", "Overtime Reason"})

			for _, v := range tableData {
				table.Append(v)
			}
			table.Render() // Send output

			return nil
		},
	}
}
