//Package quotation links up with the quotation table in the database
package quotation

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type QuoteTable struct {
	ID          string
	Quotation   string
	Customer    string
	Seller      string
	NameOfPhone string
}

//InsertQuotation insert quotation from seller in the quotation database
func InsertQuotation(db *sql.DB, customer, seller, id, quotation, nameOfPhone string) error {
	query := "INSERT INTO quotations (ID, Quotation, Customer, Seller, PhoneName) VALUES (?,?,?,?,?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, quotation, customer, seller, nameOfPhone)
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
		var q QuoteTable
		err = results.Scan(&q.ID, &q.Quotation, &q.Customer, &q.Seller, &q.NameOfPhone)
		if err != nil {
			return idList, err
		}
		idList = append(idList, q.ID)
	}
	return idList, nil
}

//GetCustomerQuote will retrieve information from quotations table where the customer field is the customer
func GetCustomerQuote(db *sql.DB, customer string) ([]QuoteTable, error) {
	list := []QuoteTable{}
	results, err := db.Query("SELECT * FROM quotations WHERE Customer = ?", customer)
	if err != nil {
		return list, err
	}
	defer results.Close()
	for results.Next() {
		var q QuoteTable
		err = results.Scan(&q.ID, &q.Quotation, &q.Customer, &q.Seller, &q.NameOfPhone)
		if err != nil {
			return list, err
		}
		list = append(list, q)
	}
	return list, nil
}

//Delete removes entry tagged to the transaction id
func Delete(db *sql.DB, id string) error {
	query := "DELETE FROM quotations WHERE ID = ?"
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
		return fmt.Errorf("ID not found")
	}
	return nil
}
