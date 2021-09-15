package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
	"time"
	"workday/internal/worklog"
)

var (
	ParamFinishTime     string
	ParamOvertimeReason string
	ParamFinishEarly    bool
)

func finishCommand() *cobra.Command {
	finishCommand := &cobra.Command{
		Use:   "finish",
		Short: "Records the workday finishing",
		RunE: func(cmd *cobra.Command, args []string) error {
			finishTime := time.Now()
			if ParamFinishTime != "" {
				finishTimeComponents := strings.Split(ParamFinishTime, ":")
				hour, err := strconv.Atoi(finishTimeComponents[0])
				if err != nil {
					return err
				}
				minute, err := strconv.Atoi(finishTimeComponents[0])
				if err != nil {
					return err
				}
				finishTime = time.Date(finishTime.Year(), finishTime.Month(), finishTime.Day(), hour, minute, 0, 0, finishTime.Location())
			}

			err := wlHandler.CompleteWorkday(finishTime, ParamOvertimeReason, ParamFinishEarly)
			if err != nil {
				if err == worklog.ErrWorkdayDurationIncomplete {
					return errors.New("you have not yet finished your full workday, use the --finishEarly parameter to use accrued overtime to finish early")
				} else if err == worklog.ErrOvertimeNotSpecified {
					return errors.New("you done overtime today, you need to specify a reason for the overtime")
				} else {
					return err
				}
			}

			fmt.Println(fmt.Sprintf("Worday finsihed at %s", finishTime.Format(time.RFC1123)))
			return nil
		},
	}

	finishCommand.Flags().StringVarP(&ParamFinishTime, "finishTime", "t", "", "hh:mm")
	finishCommand.Flags().StringVarP(&ParamOvertimeReason, "overtimeReason", "m", "", "Specifies a reason for overtime")
	finishCommand.Flags().BoolVarP(&ParamFinishEarly, "finishEarly", "e", false, "Use this parameter to finish early using accrued overtime")
	return finishCommand
}
