package secretkey

import (
	"database/sql"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

var td = secret{Type: "type1", Value: "value1"}

func TestGetKey(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()
	query := "SELECT \\* FROM secret WHERE Type = \\?"

	rows := sqlmock.NewRows([]string{"Type", "Value"}).
		AddRow(td.Type, td.Value)

	mock.ExpectQuery(query).WithArgs(td.Type).WillReturnRows(rows)

	actual, err := GetKey(db, td.Type)

	assert.Equal(t, td.Value, actual)
	assert.NoError(t, err)
}
