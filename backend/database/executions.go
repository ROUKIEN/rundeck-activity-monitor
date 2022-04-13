package database

import (
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"database/sql"
	"fmt"
	"os"

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

func Save(e *spec.Execution) error {
	connStr := "postgres://rad_user:rad_p4ss@localhost:5432/rad?sslmode=disable"
	// connStr := "user=rad_user dbname=rad_p4ss sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	id := 1
	rows, err := db.Query("SELECT id FROM rad_executions_history WHERE id > $1", id)
	if err != nil {
		return err
	}

	for rows.Next() {
		var row spec.Execution
		if err := rows.Scan(&row.ID); err != nil {
			return err
		}

		fmt.Printf("[%d]\n", row.ID)
	}

	return nil
}
