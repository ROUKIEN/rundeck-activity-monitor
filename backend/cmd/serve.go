package cmd

import (
	"ROUKIEN/rundeck-activity-monitor/database"
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

//go:embed static/*
var staticApp embed.FS

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name:   "serve",
		Usage:  "scrape rundeck instances",
		Action: serveExecute,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Value: "4000",
				Usage: "server will listen on that port",
			},
		},
	}
}

func getFiltersHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Fetching filters\n")
		query := req.URL.Query()
		begin := time.Now().Add(time.Duration(-4) * time.Hour)
		end := time.Now()

		if query.Has("begin") {
			if n, err := strconv.ParseInt(query.Get("begin"), 10, 64); err == nil {
				begin = time.Unix(n/1000, 0)
			}
		}

		if query.Has("end") {
			if n, err := strconv.ParseInt(query.Get("end"), 10, 64); err == nil {
				end = time.Unix(n/1000, 0)
			}
		}

		fmt.Printf("%s\n", begin)
		fmt.Printf("%s\n", end)
		filters, err := database.FindFilters(db, begin, end)
		if err != nil {
			internalError := http.StatusInternalServerError
			http.Error(w, err.Error(), internalError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(filters)
		}
	}
}

func getExecutionsHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Fetching executions\n")
		query := req.URL.Query()
		begin := time.Now().Add(time.Duration(-4) * time.Hour)
		end := time.Now()

		if query.Has("begin") {
			if n, err := strconv.ParseInt(query.Get("begin"), 10, 64); err == nil {
				begin = time.Unix(n/1000, 0)
			}
		}

		if query.Has("end") {
			if n, err := strconv.ParseInt(query.Get("end"), 10, 64); err == nil {
				end = time.Unix(n/1000, 0)
			}
		}

		fmt.Printf("%s\n", begin)
		fmt.Printf("%s\n", end)
		executions, err := database.FindExecutions(
			db,
			begin,
			end,
		)
		if err != nil {
			internalError := http.StatusInternalServerError
			http.Error(w, err.Error(), internalError)
		} else {
			fmt.Printf("%d executions returned\n", len(executions))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(executions)
		}
	}
}

func serveExecute(c *cli.Context) error {
	db, err := database.Db()
	if err != nil {
		return err
	}

	fsys := fs.FS(staticApp)
	contentStatic, err := fs.Sub(fsys, "static")
	if err != nil {
		return err
	}

	fs := http.FileServer(http.FS(contentStatic))
	http.Handle("/", fs)

	http.HandleFunc("/api/executions", getExecutionsHandler(db))
	http.HandleFunc("/api/filters", getFiltersHandler(db))
	port := 4000
	fmt.Printf("Listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return nil
}
