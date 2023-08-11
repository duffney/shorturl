package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/duffney/shorturl/internal/data"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port     int
	workerID int64
	dsn      string
	// db       struct {
	// 	dsn string
	// }
}

type application struct {
	config config
	logger *log.Logger
	dbMap  map[int64]Shorten
	models data.Models
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Int64Var(&cfg.workerID, "workerID", 1, "Worker ID") //Change to env var later os.Getenv("WORKER_ID")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("SHORTURL_DB_DSN"), "PostgreSql DSN")
	flag.Parse()
	fmt.Println(cfg.dsn)

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		dbMap:  make(map[int64]Shorten), // Replace with database later
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", strconv.Itoa(cfg.port)),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %d", cfg.port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
