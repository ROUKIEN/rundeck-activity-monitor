package cmd

import (
	"ROUKIEN/rundeck-activity-monitor/database"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewDatabaseCmd() *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "database operations",
		Subcommands: []*cli.Command{
			{
				Name:  "update",
				Usage: "ensure RAM database is up to date",
				Action: func(c *cli.Context) error {
					log.Info("Checking database status...")

					db, err := database.Db()
					if err != nil {
						return err
					}
					// 1. check if table already exists
					stmt, err := db.Prepare(database.SQL_RAD_TABLE_EXISTS)
					if err != nil {
						return err
					}

					defer stmt.Close()
					result := stmt.QueryRow()
					var tableExists bool

					if err := result.Scan(&tableExists); err != nil {
						return err
					}

					if tableExists {
						log.Info("Tables are already up to date.")
						return nil
					}

					createStmt, err := db.Prepare(database.SQL_RAD_TABLE_STRUCTURE)
					if err != nil {
						return err
					}
					defer createStmt.Close()

					if _, err := createStmt.Exec(); err != nil {
						return err
					}

					log.Info("Tables created.")

					return nil
				},
			},
		},
	}
}
