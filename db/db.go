package db

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const MSGS_DB = "msgs"

var DB *sql.DB

func SetConn(port, user, pass, dbname string) *sql.DB {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   "localhost:" + port,
		User:   url.UserPassword(user, pass),
		Path:   dbname,
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		log.Fatalln("Could not open DB:", err)
	}
	//defer closeConn()

	if err := db.PingContext(context.Background()); err != nil {
		log.Fatalln("Could not ping DB:", err)
	}

	return db
}

func OnExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		err := DB.Close()
		if err != nil {
			log.Println("couldnt close connection", err)
		}
		log.Println("connection was closed")
		os.Exit(1)
	}()
}

func (msgToFind *Msg) isIn(msgs []Msg) bool {
	for _, m := range msgs {
		if m.Id == msgToFind.Id {
			return true
		}
	}
	return false
}

func (msgToCheck *Msg) timestampIs(t string) bool {
	return strings.Contains(msgToCheck.Timestamp, t)
}
