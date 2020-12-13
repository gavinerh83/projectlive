//Package users interacts with the database to insert and retrieve user information
package users

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//User contains fields for user signups
type User struct {
	Username  string
	Password  string
	Company   string
	IsCompany string
}

//InsertRecord inserts new user record signups into the database
func InsertRecord(db *sql.DB, username, password, company, iscompany string) error {
	query := fmt.Sprintf("INSERT INTO users (Username, Password, Company, IsCompany) VALUES ('%s', '%s', '%s', '%s')", username, password, company, iscompany)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

//GetRecord retrieves information from the database
func GetRecord(db *sql.DB) (map[string]User, error) {
	var results *sql.Rows
	var err error
	userMap := map[string]User{}
	results, err = db.Query("SELECT * FROM users")
	if err != nil {
		return userMap, err
	}
	defer results.Close()
	for results.Next() {
		var users User
		err = results.Scan(&users.Username, &users.Password, &users.Company, &users.IsCompany)
		if err != nil {
			return userMap, err
		}
		userMap[users.Username] = users
		if err != nil {
			return userMap, err
		}
	}
	return userMap, nil
}

//DeleteRecord deletes record from the database
func DeleteRecord(db *sql.DB, username string) error {
	query := fmt.Sprintf("DELETE FROM users WHERE Username ='%s'", username)
	result, err := db.Exec(query) //try exec and get the resuults
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("Username not found")
	}
	return nil
}

//SearchRecord search for record in the database and returns the user information
func SearchRecord(db *sql.DB, username string) (User, error) {
	var results *sql.Rows
	var users User
	var err error
	results, err = db.Query("SELECT * FROM users WHERE Username = ?", username)
	if err != nil {
		return users, err
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&users.Username, &users.Password, &users.Company, &users.IsCompany)
		if err != nil {
			return users, err
		}
	}
	return users, nil
}

//UpdateRecord updates the password of the user after resetting
func UpdateRecord(db *sql.DB, password, username string) error {
	query2 := fmt.Sprintf("UPDATE users SET Password ='%s' WHERE Username ='%s'", password, username)
	_, err := db.Query(query2)
	if err != nil {
		return err
	}
	return nil
}
