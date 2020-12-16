package users

import (
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

type data struct {
	Username  string
	Password  string
	Company   string
	IsCompany string
}

var td = data{"jeff", "password2", "nil", "false"}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestInsertRecord(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	query := "INSERT INTO users \\(Username, Password, Company, IsCompany\\)  VALUES \\(\\?, \\?, \\?, \\?\\)"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(td.Username, td.Password, td.Company, td.IsCompany).WillReturnResult(sqlmock.NewResult(0, 1))
	err := InsertRecord(db, td.Username, td.Password, td.Company, td.IsCompany)
	assert.NoError(t, err)
}

func TestGetRecord(t *testing.T) {
	var expectedm = map[string]User{}
	db, mock := NewMock()
	defer db.Close()
	query := "SELECT \\* FROM users"

	rows := sqlmock.NewRows([]string{"Username", "Password", "Company", "IsCompany"}).
		AddRow(td.Username, td.Password, td.Company, td.IsCompany)

	mock.ExpectQuery(query).WillReturnRows(rows)
	m, err := GetRecord(db)
	expectedm[td.Username] = User{td.Username, td.Password, td.Company, td.IsCompany}

	assert.Equal(t, expectedm, m)
	assert.NoError(t, err)
}

func TestSearchRecord(t *testing.T) {
	//create new db
	db, mock := NewMock()
	defer db.Close()
	//create rows
	//create query
	query := "SELECT \\* FROM users WHERE Username = \\?"
	//create rows in db
	rows := sqlmock.NewRows([]string{"Username", "Password", "Company", "IsCompany"}).
		AddRow(td.Username, td.Password, td.Company, td.IsCompany)

	mock.ExpectQuery(query).WithArgs(td.Username).WillReturnRows(rows)

	//actual result
	user, err := SearchRecord(db, td.Username)
	//construct expected result
	var expecteduser User
	expecteduser = User{td.Username, td.Password, td.Company, td.IsCompany}
	assert.Equal(t, expecteduser, user)
	assert.NoError(t, err)
}
func TestDeleteRecord(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	query := "DELETE FROM users WHERE Username = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(td.Username).WillReturnResult(sqlmock.NewResult(0, 1))

	err := DeleteRecord(db, td.Username)
	assert.NoError(t, err)
}

func TestUpdateRecord(t *testing.T) {
	//create new db
	db, mock := NewMock()
	//create query
	query := "UPDATE users SET Password = \\? WHERE Username = \\?"

	//prepare statement
	prep := mock.ExpectPrepare(query)
	//execute
	prep.ExpectExec().WithArgs(td.Password, td.Username).WillReturnResult(sqlmock.NewResult(0, 1))
	err := UpdateRecord(db, td.Password, td.Username)
	assert.NoError(t, err)
}
