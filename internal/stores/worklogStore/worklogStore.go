package worklogStore

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"workday/internal/worklog"
)

const database = "sqlite3"

func NewWorkLogStore() (worklog.Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	db, err := setupDatabase(path.Join(homeDir, ".workday.db"))
	if err != nil {
		return nil, err
	}

	return workLogStore{
		db: db,
	}, nil
}

func setupDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open(database, dbPath)
	if err != nil {
		return nil, err
	}
	if err := performMigrations(db); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}
	return db, nil
}

func performMigrations(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, new(sqlite3.Config))
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/stores/worklogStore/migrations", database, driver)
	return m.Up()
}

type workLogStore struct {
	db *sql.DB
}

func (w workLogStore) GetCurrentLogEntry() (*worklog.Entry, error) {
	rows, err := w.db.Query("SELECT wl.ROWID, wl.start_time FROM work_log wl WHERE wl.finish_time IS NULL")
	if err != nil {
		return nil, err
	}

	var currentLogEntry *worklog.Entry = nil
	for rows.Next() {
		currentLogEntry = new(worklog.Entry)
		if err := rows.Scan(&(*currentLogEntry).ID, &(*currentLogEntry).StartTime); err != nil {
			return nil, err
		}
	}
	return currentLogEntry, nil
}

func (w workLogStore) ReadCompletedWorkLog() (*worklog.WorkLog, error) {
	rows, err := w.db.Query("SELECT wl.ROWID, wl.start_time, wl.finish_time, wl.overtime_reason FROM work_log wl WHERE wl.finish_time IS NOT NULL ORDER BY wl.finish_time")
	if err != nil {
		return nil, err
	}

	var log *worklog.WorkLog = nil
	for rows.Next() {
		if log == nil {
			log = new(worklog.WorkLog)
		}
		var workEntry worklog.Entry
		if err := rows.Scan(&workEntry.ID, &workEntry.StartTime, &workEntry.FinishTime, &workEntry.OvertimeReason); err != nil {
			return nil, err
		}
		*log = append(*log, workEntry)
	}
	return log, nil
}

func (w workLogStore) AddNewWorkLogEntry(entry *worklog.Entry) error {
	_, err := w.db.Exec("INSERT INTO work_log (start_time) values (?)", entry.StartTime)
	return err
}

func (w workLogStore) UpdateWorkLogEntry(entry *worklog.Entry) error {
	_, err := w.db.Exec("UPDATE work_log SET finish_time = ?, overtime_reason = ? WHERE ROWID = ?", entry.FinishTime, entry.OvertimeReason, entry.ID)
	return err
}
