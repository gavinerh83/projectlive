package submissions

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Condition struct {
	Customer            string
	Name                string
	Storage             string
	Housing             string
	Screen              string
	OriginalAccessories string
	OtherIssues         string
	ID                  string
}

//InsertDetails inserts phone information into database
func InsertDetails(db *sql.DB, customer, name, storage, housing, screen, originalaccessories, otherissues, ID string) error {
	query := "INSERT INTO submissions (Username, Name, Storage, Housing, Screen, Original_Accessories, Other_Issues, ID) VALUES (?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(customer, name, storage, housing, screen, originalaccessories, otherissues, ID)
	if err != nil {
		return err
	}
	return nil
}

//GetDetails retrieves information from the database
func GetDetails(db *sql.DB) (map[string]Condition, error) {
	submissions := map[string]Condition{}
	results, err := db.Query("SELECT * FROM submissions")
	if err != nil {
		return submissions, err
	}
	defer results.Close()
	for results.Next() {
		var phone Condition
		err = results.Scan(&phone.ID, &phone.Customer, &phone.Name, &phone.Storage, &phone.Housing, &phone.Screen, &phone.OriginalAccessories, &phone.OtherIssues)
		if err != nil {
			return submissions, err
		}
		submissions[phone.ID] = phone
	}
	return submissions, nil
}

//GetID uses the transaction id and return specific phone information
func GetID(db *sql.DB, id string) (Condition, error) {
	var phone Condition
	err := db.QueryRow("SELECT * FROM submissions WHERE ID = ?", id).Scan(&phone.ID, &phone.Customer, &phone.Name, &phone.Storage, &phone.Housing, &phone.Screen, &phone.OriginalAccessories, &phone.OtherIssues)
	if err != nil {
		return phone, err
	}
	return phone, nil
}

//Delete removes the row with the id specified
func Delete(db *sql.DB, id string) error {
	query := "DELETE FROM submissions WHERE ID = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
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
