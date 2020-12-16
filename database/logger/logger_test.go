package logger

import (
	"database/sql"
	"log"
	"testing"
	"time"

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

func TestLogging(t *testing.T) {
	db, mock := NewMock()
	tn := time.Now().String()
	tn = tn[:20]
	msg := "there is error"
	query := "INSERT INTO logging \\(DateTime, Message\\) VALUES \\(\\?,\\?\\)"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(tn, msg).WillReturnResult(sqlmock.NewResult(0, 1))

	err := Logging(db, msg)
	assert.NoError(t, err)
}
