package quotation

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

var td = QuoteTable{ID: "id1", Quotation: "quotation1", Customer: "customer1", Seller: "seller1", NameOfPhone: "phone1"}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestInsertQuotation(t *testing.T) {
	db, mock := NewMock()
	query := "INSERT INTO quotations \\(ID, Quotation, Customer, Seller, PhoneName\\) VALUES \\(\\?,\\?,\\?,\\?,\\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(td.ID, td.Quotation, td.Customer, td.Seller, td.NameOfPhone).WillReturnResult(sqlmock.NewResult(0, 1))

	err := InsertQuotation(db, td.Customer, td.Seller, td.ID, td.Quotation, td.NameOfPhone)
	assert.NoError(t, err)
}

func TestSearchSeller(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT \\* FROM quotations WHERE Seller = \\?"

	rows := sqlmock.NewRows([]string{"ID", "Quotation", "Customer", "Seller", "PhoneName"}).
		AddRow(td.ID, td.Quotation, td.Customer, td.Seller, td.NameOfPhone)

	mock.ExpectQuery(query).WithArgs(td.Seller).WillReturnRows(rows)

	actual, err := SearchSeller(db, td.Seller)
	var expected []string
	expected = append(expected, td.ID)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestGetCustomerQuote(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT \\* FROM quotations WHERE Customer = \\?"
	rows := sqlmock.NewRows([]string{"ID", "Quotation", "Customer", "Seller", "PhoneName"}).
		AddRow(td.ID, td.Quotation, td.Customer, td.Seller, td.NameOfPhone)

	mock.ExpectQuery(query).WithArgs(td.Customer).WillReturnRows(rows)

	actual, err := GetCustomerQuote(db, td.Customer)

	var expected []QuoteTable
	expected = append(expected, td)
	assert.Equal(t, expected, actual)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	db, mock := NewMock()
	query := "DELETE FROM quotations WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(td.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := Delete(db, td.ID)
	assert.NoError(t, err)
}
