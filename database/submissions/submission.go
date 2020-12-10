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
	query := fmt.Sprintf("INSERT INTO submissions (Username, Name, Storage, Housing, Screen, Original_Accessories, Other_Issues, ID) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", customer, name, storage, housing, screen, originalaccessories, otherissues, ID)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

//GetRecord retrieves information from the database
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
	results, err := db.Query("SELECT * FROM submissions WHERE ID = ?", id)
	if err != nil {
		return phone, err
	}
	defer results.Close()
	for results.Next() {
		err = results.Scan(&phone.ID, &phone.Customer, &phone.Name, &phone.Storage, &phone.Housing, &phone.Screen, &phone.OriginalAccessories, &phone.OtherIssues)
		if err != nil {
			return phone, err
		}
	}
	return phone, nil
}

// //DeleteRecord deletes record from the submissions table and insert into pastsubmissions table
// func transferDetails(db *sql.DB, customerName, sellername string) error {
// 	//retrieve record
// 	var results *sql.Rows
// 	var err error
// 	submissions := map[string]Condition{}
// 	results, err = db.Query("SELECT * FROM submissions WHERE Customer = ?", customerName)
// 	if err != nil {
// 		return err
// 	}
// 	defer results.Close()
// 	var phone Condition
// 	for results.Next() {
// 		err = results.Scan(&phone.Customer, &phone.Name, &phone.Storage, &phone.Housing, &phone.OriginalAccessories, &phone.OtherIssues)
// 		if err != nil {
// 			return err
// 		}
// 		submissions[phone.Customer] = phone
// 	}
// 	//Delete data from submissions database
// 	query := fmt.Sprintf("DELETE FROM submissions WHERE Customer ='%s'", customerName)
// 	result, err := db.Exec(query) //try exec and get the resuults
// 	if err != nil {
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("Customer not found")
// 	}
// 	//Insert into pastsubmissions database
// 	query = fmt.Sprintf("INSERT INTO pastsubmissions (Customer, Seller, Name, Storage, Housing, Original_Accessories, Other_Issues) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')", customerName, sellername, phone., housing, originalaccessories, otherissues)
// 	_, err = db.Query(query)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// 	return nil
// }
