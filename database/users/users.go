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
	query := "INSERT INTO users (Username, Password, Company, IsCompany) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, password, company, iscompany)
	if err != nil {
		return err
	}
	return nil
}

//GetRecord retrieves information from the database
func GetRecord(db *sql.DB) (map[string]User, error) {
	var rows *sql.Rows
	var err error
	userMap := map[string]User{}
	rows, err = db.Query("SELECT * FROM users")
	if err != nil {
		return userMap, err
	}
	defer rows.Close()
	for rows.Next() {
		var users User
		err = rows.Scan(&users.Username, &users.Password, &users.Company, &users.IsCompany)
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
	query := "DELETE FROM users WHERE Username = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(username)
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
	var users User
	var err error
	err = db.QueryRow("SELECT * FROM users WHERE Username = ?", username).Scan(&users.Username, &users.Password, &users.Company, &users.IsCompany)
	if err != nil {
		return users, err
	}
	return users, nil
}

//UpdateRecord updates the password of the user after resetting
func UpdateRecord(db *sql.DB, password, username string) error {
	query2 := "UPDATE users SET Password = ? WHERE Username = ?"
	stmt, err := db.Prepare(query2)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(password, username)
	if err != nil {
		return err
	}
	return nil
}

//RetrieveSeller retrieves seller information
func RetrieveSeller(db *sql.DB) ([]User, error) {
	var users []User
	b := "true"
	rows, err := db.Query("SELECT * FROM users WHERE IsCompany = ?", b)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Username, &user.Password, &user.Company, &user.IsCompany)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
