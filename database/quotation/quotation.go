package quotation

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func InsertQuotation(db *sql.DB, customer, seller, id string, quotation int) error {
	query := fmt.Sprintf("INSERT INTO submissions (ID, Quotation, Customer, Seller) VALUES ('%s', %d, '%s', '%s')", id, quotation, customer, seller)
	_, err := db.Query(query)
	if err != nil {
		return err
	}
	return nil
}
