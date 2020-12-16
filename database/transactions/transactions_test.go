package transactions

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var sub = PSubmissions{ID: "1", Customer: "customer1", Seller: "seller1", PhoneName: "phonename1", Storage: "storage1", Housing: "housing1", Screen: "screen1", OriginalAccessories: "OriginalAccessories1", OtherIssues: "otherissues1", Quotation: "quotation1", DateTime: "date1"}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
func TestInsertTransaction(t *testing.T) {
	db, mock := NewMock()
	//create query
	query := "INSERT INTO pastsubmissions \\(ID, Customer, Seller, Name, Storage, Housing, Screen, Original_Accessories, Other_Issues, Quotation, DateTime\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?,\\?\\)"
	//prepare
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(sub.ID, sub.Customer, sub.Seller, sub.PhoneName, sub.Storage, sub.Housing, sub.Screen, sub.OriginalAccessories, sub.OtherIssues, sub.Quotation, sub.DateTime).WillReturnResult(sqlmock.NewResult(0, 1))
	//exec the actual
	err := InsertTransaction(db, sub.ID, sub.Customer, sub.Seller, sub.PhoneName, sub.Storage, sub.Housing, sub.Screen, sub.OriginalAccessories, sub.OtherIssues, sub.Quotation, sub.DateTime)
	assert.NoError(t, err)
}

func TestGetCustomer(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT \\* FROM pastsubmissions WHERE Customer = \\?"
	//for selecting from database, need to first create rows
	rows := sqlmock.NewRows([]string{"ID", "Customer", "Seller", "PhoneName", "Storage", "Housing", "Screen", "Original_Accessories", "Other_Issues", "Quotation", "DateTime"}).
		AddRow(sub.ID, sub.Customer, sub.Seller, sub.PhoneName, sub.Storage, sub.Housing, sub.Screen, sub.OriginalAccessories, sub.OtherIssues, sub.Quotation, sub.DateTime)

	mock.ExpectQuery(query).WithArgs(sub.Customer).WillReturnRows(rows)

	//actual result
	ss, err := GetCustomer(db, sub.Customer)

	//create expected result
	var expected []PSubmissions
	expected = append(expected, sub)
	assert.Equal(t, expected, ss)
	assert.NoError(t, err)
}

func TestGetSeller(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT \\* FROM pastsubmissions WHERE Seller = \\?"
	//for selecting from database, need to first create rows
	rows := sqlmock.NewRows([]string{"ID", "Customer", "Seller", "PhoneName", "Storage", "Housing", "Screen", "Original_Accessories", "Other_Issues", "Quotation", "DateTime"}).
		AddRow(sub.ID, sub.Customer, sub.Seller, sub.PhoneName, sub.Storage, sub.Housing, sub.Screen, sub.OriginalAccessories, sub.OtherIssues, sub.Quotation, sub.DateTime)

	mock.ExpectQuery(query).WithArgs(sub.Seller).WillReturnRows(rows)

	//actual result
	ss, err := GetSeller(db, sub.Seller)

	//create expected result
	var expected []PSubmissions
	expected = append(expected, sub)
	assert.Equal(t, expected, ss)
	assert.NoError(t, err)
}
