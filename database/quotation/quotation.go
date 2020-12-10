//Package quotation links up with the quotation table in the database
package quotation

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type quoteTable struct {
	ID        string
	Quotation string
	Customer  string
	Seller    string
}

//InsertQuotation insert quotation from seller in the quotation database
func InsertQuotation(db *sql.DB, customer, seller, id, quotation string) error {
	query := fmt.Sprintf("INSERT INTO quotations (ID, Quotation, Customer, Seller) VALUES ('%s', '%s', '%s', '%s')", id, quotation, customer, seller)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

//SearchSeller searches the database if the seller name exist, if exist, seller will not see that order.
//SearchSeller returns the transaction id
func SearchSeller(db *sql.DB, seller string) ([]string, error) {
	idList := []string{}
	results, err := db.Query("SELECT * FROM quotations WHERE Seller = ?", seller)
	if err != nil {
		return idList, err
	}
	defer results.Close()
	for results.Next() {
		var q quoteTable
		err = results.Scan(&q.ID, &q.Quotation, &q.Customer, &q.Seller)
		if err != nil {
			return idList, err
		}
		idList = append(idList, q.ID)
	}
	return idList, nil
}
