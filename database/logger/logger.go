//Package logger is a program that logs custom error messages in a database.
package logger

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Logging takes in a error message and prints the output to a postgres database
func Logging(db *sql.DB, msg string) error {
	t := time.Now().String()
	t = t[:28]
	query := fmt.Sprintf("INSERT INTO logging (DateTime, Message) VALUES ('%s', '%s')", t, msg)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}
