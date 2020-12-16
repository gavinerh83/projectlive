//Package logger is a program that logs custom error messages in a database.
package logger

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Logging takes in a error message and prints the output to a postgres database
func Logging(db *sql.DB, msg string) error {
	t := time.Now().String()
	t = t[:20]
	query := "INSERT INTO logging (DateTime, Message) VALUES (?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t, msg)
	if err != nil {
		return err
	}
	return nil
}
