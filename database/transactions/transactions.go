package transactions

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//PSubmissions contains the fields of the past transactions
type PSubmissions struct {
	ID                  string
	Customer            string
	Seller              string
	PhoneName           string
	Storage             string
	Housing             string
	Screen              string
	OriginalAccessories string
	OtherIssues         string
	Quotation           string
	DateTime            string
}

//InsertTransaction inserts successful transactions into the pastSubmission table
func InsertTransaction(db *sql.DB, id, customer, seller, name, storage, housing, screen, originalAccessories, otherissues, quotation, datetime string) error {
	query := "INSERT INTO pastsubmissions (ID, Customer, Seller, Name, Storage, Housing, Screen, Original_Accessories, Other_Issues, Quotation, DateTime) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, customer, seller, name, storage, housing, screen, originalAccessories, otherissues, quotation, datetime)
	if err != nil {
		return err
	}
	return nil
}

//GetCustomer accepts username and returns the transaction information
func GetCustomer(db *sql.DB, username string) ([]PSubmissions, error) {
	transacs := []PSubmissions{}
	results, err := db.Query("SELECT * FROM pastsubmissions WHERE Customer = ?", username)
	if err != nil {
		return transacs, err
	}
	defer results.Close()
	for results.Next() {
		var transac PSubmissions
		err = results.Scan(&transac.ID, &transac.Customer, &transac.Seller, &transac.PhoneName, &transac.Storage, &transac.Housing, &transac.Screen, &transac.OriginalAccessories, &transac.OtherIssues, &transac.Quotation, &transac.DateTime)
		if err != nil {
			return transacs, err
		}
		transacs = append(transacs, transac)
	}
	return transacs, nil
}

//GetSeller accepts username and returns the transaction information
func GetSeller(db *sql.DB, username string) ([]PSubmissions, error) {
	transacs := []PSubmissions{}
	results, err := db.Query("SELECT * FROM pastsubmissions WHERE Seller = ?", username)
	if err != nil {
		return transacs, err
	}
	defer results.Close()
	for results.Next() {
		var transac PSubmissions
		err = results.Scan(&transac.ID, &transac.Customer, &transac.Seller, &transac.PhoneName, &transac.Storage, &transac.Housing, &transac.Screen, &transac.OriginalAccessories, &transac.OtherIssues, &transac.Quotation, &transac.DateTime)
		if err != nil {
			return transacs, err
		}
		transacs = append(transacs, transac)
	}
	return transacs, nil
}
