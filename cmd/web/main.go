package main

import (
	"flag"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log/slog"
	"net/http"
	"os"
	"snippetbox.zacharyad/internal/data"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	// this is just a pointer value to the passed in OR default value for the port
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:abc123@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Starting server", slog.String("addr", ":4000"))

	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
