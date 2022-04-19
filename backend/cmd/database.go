package cmd

import (
	"ROUKIEN/rundeck-activity-monitor/database"
	"fmt"

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
					fmt.Printf("Checking database status...\n")

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
						fmt.Printf("Tables are already up to date.\n")
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

					fmt.Printf("Tables created.\n")

					return nil
				},
			},
		},
	}
}
