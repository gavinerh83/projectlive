package oauth

import (
	"database/sql"
	"fmt"
	"time"
)

type authentication struct {
	Time   string
	TempID string
}

//InsertTempID inserts temporary id into database
func InsertTempID(db *sql.DB, id string) error {
	t := time.Now().String()
	t = t[:28]
	query := fmt.Sprintf("INSERT INTO oauth2 (Time, TempID) VALUES ('%s', '%s')", t, id)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

//GetTempID returns the temp id previously stored
func GetTempID(db *sql.DB) (map[string]string, error) {
	var m = map[string]string{}
	results, err := db.Query("SELECT * FROM oauth2")
	if err != nil {
		return m, err
	}
	defer results.Close()
	var k authentication
	for results.Next() {
		err = results.Scan(&k.Time, &k.TempID)
		if err != nil {
			return m, err
		}
		m[k.TempID] = k.Time
	}
	return m, nil
}
