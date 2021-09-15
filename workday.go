package main

import (
	"fmt"
	"os"
	"workday/cmd"
	"workday/internal/stores/worklogStore"
	"workday/internal/worklog"
)

func main() {
	workLogStore, err := worklogStore.NewWorkLogStore()
	if err != nil {
		handleErr(err)
	}

	if err := cmd.Execute(worklog.NewHandler(workLogStore, worklog.WorkdayConfig{WorkdayDuration: 8})); err != nil {
		handleErr(err)
	}
}

func handleErr(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}