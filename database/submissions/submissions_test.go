package submissions

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
)

var c = Condition{Customer: "customer1", Name: "name1", Storage: "storage1", Housing: "housing", Screen: "screen1", OriginalAccessories: "origin", OtherIssues: "other1", ID: "id1"}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGetDetails(t *testing.T) {
	db, mock := NewMock()
	//create query
	query := "SELECT \\* FROM submissions"
	//create rows
	rows := sqlmock.NewRows([]string{"ID", "Username", "Name", "Storage", "Housing", "Screen", "Original_Accessories", "Other_Issues"}).
		AddRow(c.ID, c.Customer, c.Name, c.Storage, c.Housing, c.Screen, c.OriginalAccessories, c.OtherIssues)

	mock.ExpectQuery(query).WillReturnRows(rows)
	//actual
	m, err := GetDetails(db)
	var expectedm = map[string]Condition{}
	expectedm[c.ID] = c
	assert.Equal(t, expectedm, m)
	assert.NoError(t, err)
}

func TestGetID(t *testing.T) {
	db, mock := NewMock()
	query := "SELECT \\* FROM submissions WHERE ID = \\?"
	rows := sqlmock.NewRows([]string{"ID", "Username", "Name", "Storage", "Housing", "Screen", "Original_Accessories", "Other_Issues"}).
		AddRow(c.ID, c.Customer, c.Name, c.Storage, c.Housing, c.Screen, c.OriginalAccessories, c.OtherIssues)

	mock.ExpectQuery(query).WithArgs(c.ID).WillReturnRows(rows)

	actual, err := GetID(db, c.ID)
	assert.NoError(t, err)
	assert.Equal(t, c, actual)
}

func TestDelete(t *testing.T) {
	db, mock := NewMock()
	query := "DELETE FROM submissions WHERE ID = \\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(c.ID).WillReturnResult(sqlmock.NewResult(0, 1))

	err := Delete(db, c.ID)
	assert.NoError(t, err)
}
