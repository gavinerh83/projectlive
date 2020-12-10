package transactions

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

//InsertTransaction inserts successful transactions into the pastSubmission table
func InsertTransaction(db *sql.DB, id, customer, seller, name, storage, housing, screen, originalAccessories, otherissues, datetime string) error {
	query := fmt.Sprintf("INSERT INTO pastsubmissions (ID, Customer, Seller, Name, Storage, Housing, Screen, Original_Accessories, Other_Issues, DateTime) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')", id, customer, seller, name, storage, housing, screen, originalAccessories, otherissues, datetime)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}
