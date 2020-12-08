//Package logger is a program that logs custom error messages in a database.
package logger

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

var (
	//Error is a type log
	Error *log.Logger
	mu    sync.Mutex
)

//Logging takes in a error message and prints the output to a postgres database
func Logging(msg string) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	mu.Lock()
	db, err := sql.Open("postgres", "user=postgres password=password host=127.0.0.1 port=5432 dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("The connection to the DB was successfully initialised!")
	}
	connectivity := db.Ping()
	if connectivity != nil {
		panic(err)
	} else {
		fmt.Println("good to go")
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hashlog(
		id SERIAL PRIMARY KEY NOT NULL,
		content text
		)`)
	if err != nil {
		panic("Error in creating table")
	}
	stmt, err := db.Prepare(`INSERT INTO hashlog (content) VALUES ($1)`)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	Error = log.New(buf, "", log.Ldate|log.Ltime)
	Error.SetOutput(buf)
	Error.Println(buf, msg)
	m := buf.String()
	truncated := strings.TrimSpace(m)
	_, err = stmt.Exec(truncated)
	if err != nil {
		panic(err)
	}
	mu.Unlock()
}
