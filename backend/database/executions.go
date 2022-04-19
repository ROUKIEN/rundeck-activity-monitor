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
	ID                 int       `json:"id"`
	RundeckInstance    string    `json:"rundeck_instance"`
	RundeckProject     string    `json:"rundeck_project"`
	RundeckJob         string    `json:"rundeck_job"`
	RundeckJobId       string    `json:"rundeck_job_id"`
	ExecutionStatus    string    `json:"execution_status"`
	ExecutionId        int       `json:"execution_id"`
	ExecutionStart     time.Time `json:"execution_start"`
	ExecutionEnd       time.Time `json:"execution_end"`
	ExecutionPermalink string    `json:"execution_permalink"`
}

type SearchFilters struct {
	Instances []string `json:"instances"`
	Statuses  []string `json:"statuses"`
	Projects  []string `json:"projects"`
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

func FindExecutions(db *sql.DB, start time.Time, end time.Time) ([]*RecordedExecution, error) {

	stmt, err := db.Prepare(SQL_FIND_EXECUTIONS_BY_DATE)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(start, end)
	if err != nil {
		return nil, err
	}

	executions := make([]*RecordedExecution, 0)

	for rows.Next() {
		var executionInDb RecordedExecution
		if err := rows.Scan(
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
			return nil, err
		}

		executions = append(executions, &executionInDb)
	}

	return executions, nil
}

func FindFilters(db *sql.DB, start time.Time, end time.Time) (*SearchFilters, error) {
	instances, err := findInstancesOnTimerange(db, start, end)
	if err != nil {
		return nil, err
	}

	projects, err := findProjectsOnTimerange(db, start, end)
	if err != nil {
		return nil, err
	}

	statuses, err := findStatusesOnTimerange(db, start, end)
	if err != nil {
		return nil, err
	}

	return &SearchFilters{
		Instances: instances,
		Projects:  projects,
		Statuses:  statuses,
	}, nil
}

func findInstancesOnTimerange(db *sql.DB, start time.Time, end time.Time) ([]string, error) {
	stmt, err := db.Prepare(SQL_FIND_INSTANCES_ON_TIMERANGE)
	if err != nil {
		return nil, err
	}
	instances := make([]string, 0)
	rows, err := stmt.Query(start, end)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var instance string
		if err := rows.Scan(&instance); err == nil {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

func findProjectsOnTimerange(db *sql.DB, start time.Time, end time.Time) ([]string, error) {
	stmt, err := db.Prepare(SQL_FIND_PROJECTS_ON_TIMERANGE)
	if err != nil {
		return nil, err
	}
	projects := make([]string, 0)
	rows, err := stmt.Query(start, end)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project string
		if err := rows.Scan(&project); err == nil {
			projects = append(projects, project)
		}
	}

	return projects, nil
}

func findStatusesOnTimerange(db *sql.DB, start time.Time, end time.Time) ([]string, error) {
	stmt, err := db.Prepare(SQL_FIND_STATUSES_ON_TIMERANGE)
	if err != nil {
		return nil, err
	}
	statuses := make([]string, 0)
	rows, err := stmt.Query(start, end)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var status string
		if err := rows.Scan(&status); err == nil {
			statuses = append(statuses, status)
		}
	}

	return statuses, nil
}
