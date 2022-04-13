package database

import (
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func ConnectionString() string {
	return os.Getenv("RAM_DB_DSN")
}

func Db() (*sql.DB, error) {
	db, err := sql.Open("postgres", ConnectionString())
	if err != nil {
		return nil, err
	}

	return db, nil
}

type RecordedExecution struct {
	ID                 int
	RundeckInstance    string
	RundeckProject     string
	RundeckJob         string
	RundeckJobId       string
	ExecutionStatus    string
	ExecutionId        int
	ExecutionStart     time.Time
	ExecutionEnd       time.Time
	ExecutionPermalink string
}

func FindExecution(db *sql.DB, instance_name string, e *spec.Execution) (*RecordedExecution, error) {
	var executionInDb RecordedExecution
	stmt, err := db.Prepare(SQL_FIND_EXECUTION)
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(instance_name, e.Job.ID, e.ID)

	if err := row.Scan(
		&executionInDb.ID,
		&executionInDb.RundeckInstance,
		&executionInDb.RundeckProject,
		&executionInDb.RundeckJob,
		&executionInDb.RundeckJobId,
		&executionInDb.ExecutionStatus,
		&executionInDb.ExecutionId,
		&executionInDb.ExecutionStart,
		&executionInDb.ExecutionEnd,
		&executionInDb.ExecutionPermalink); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &executionInDb, nil
}

func SaveExecution(db *sql.DB, instance_name string, e *spec.Execution) error {
	stmt, err := db.Prepare(SQL_INSERT_EXECUTION)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(
		instance_name,
		e.Project,
		e.Job.Name,
		e.Job.ID,
		e.Status,
		e.ID,
		e.DateStarted.Date,
		e.DateEnded.Date,
		e.Permalink); err != nil {
		return err
	}

	return nil
}
