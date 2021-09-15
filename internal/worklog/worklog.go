package worklog

import (
	"errors"
	"time"
	"workday/pkg/timeHelper"
)

var (
	ErrPreviousWorkdayIncomplete = errors.New("previous workday incomplete")
	ErrWorkdayNotStarted         = errors.New("workday not yet started")
	ErrWorkdayAlreadyStarted     = errors.New("workday already started")
	ErrOvertimeNotSpecified      = errors.New("overtime reason not specified")
	ErrWorkdayDurationIncomplete = errors.New("day duration not completed")
)

type WorkdayConfig struct {
	WorkdayDuration int
}

type WorkLog []Entry

type Entry struct {
	ID             int64
	StartTime      time.Time
	FinishTime     time.Time
	OvertimeReason string
}

type Status struct {
	CurrentWorkday        *Entry
	UnusedOvertimeEntries *WorkLog
}

func (s *Status) CalculateUnusedOvertime(workdayDuration time.Duration) time.Duration {
	var overtimeDuration time.Duration
	for _, entry := range *s.UnusedOvertimeEntries {
		overtimeDuration += entry.FinishTime.Sub(entry.StartTime)
	}
	return overtimeDuration
}

type Store interface {
	GetCurrentLogEntry() (*Entry, error)
	ReadCompletedWorkLog() (*WorkLog, error)
	AddNewWorkLogEntry(entry *Entry) error
	UpdateWorkLogEntry(entry *Entry) error
}

func NewHandler(store Store, wdConfig WorkdayConfig) Handler {
	return Handler{
		store:  store,
		config: wdConfig,
	}
}

type Handler struct {
	store  Store
	config WorkdayConfig
}

func (h *Handler) GetCurrentWorkday() (*Entry, error) {
	return h.store.GetCurrentLogEntry()
}

func (h *Handler) CompleteWorkday(finishTime time.Time, overtimeReason string, useOvertimeToFinishEarly bool) error {
	currentEntry, err := h.GetCurrentWorkday()
	if err != nil {
		return err
	}
	if currentEntry == nil {
		return ErrWorkdayNotStarted
	}
	currentEntry.FinishTime = finishTime
	currentEntry.OvertimeReason = overtimeReason

	if err := h.validateWorkdayCanComplete(currentEntry, useOvertimeToFinishEarly); err != nil {
		return err
	}

	return h.store.UpdateWorkLogEntry(currentEntry)
}

func (h *Handler) StartNewWorkday() (time.Time, error) {
	entry := &Entry{
		StartTime: time.Now(),
	}

	if ce, err := h.GetCurrentWorkday(); err != nil {
		return entry.StartTime, err
	} else if ce != nil {
		if timeHelper.TimeIsDayBeforeTime(ce.StartTime, entry.StartTime) {
			return entry.StartTime, ErrPreviousWorkdayIncomplete
		} else {
			return entry.StartTime, ErrWorkdayAlreadyStarted
		}
	}

	return entry.StartTime, h.store.AddNewWorkLogEntry(entry)
}

func (h *Handler) validateWorkdayCanComplete(finishEntry *Entry, useOvertimeToFinishEarly bool) error {
	if !useOvertimeToFinishEarly && int(finishEntry.FinishTime.Sub(finishEntry.StartTime).Hours()) < h.config.WorkdayDuration {
		return ErrWorkdayDurationIncomplete
	}

	if finishEntry.FinishTime.After(finishEntry.StartTime.Add(time.Hour*time.Duration(h.config.WorkdayDuration))) && finishEntry.OvertimeReason == "" {
		return ErrOvertimeNotSpecified
	}

	return nil
}

func (h *Handler) GetStatus() (*Status, error) {
	status := new(Status)
	logs, err := h.store.ReadCompletedWorkLog()
	if err != nil {
		return nil, err
	}
	workingCount := 0
	for _, log := range *logs {
		if log.StartTime.Add(time.Hour*time.Duration(h.config.WorkdayDuration)).Sub(log.FinishTime) > 0
	}

	currentDay, err := h.store.GetCurrentLogEntry()
	if err != nil {
		return nil, err
	}

	status.UnusedOvertimeEntries = logs
	status.CurrentWorkday = currentDay
	return status, nil
}