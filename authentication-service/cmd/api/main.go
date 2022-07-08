package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "8083"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting authentication service")
	// TODO Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("cant connect to postgress")
	}

	// setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	// dsn := os.Getenv("DSN")
	dsn := "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not ready")
			counts++
		} else {
			log.Println("connected to postgress")
			return connection
		}
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("backing off for 2 secnods")
		time.Sleep(2 * time.Second)
		continue
	}
}
