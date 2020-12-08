package users

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type data struct {
	Username  string
	Password  string
	Company   string
	IsCompany string
}

var sqluser = "root"
var sqlpassword = "password"

func connectDB() *sql.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:8888)/store", sqluser, sqlpassword)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println("Good to go")
	}
	return db
}

func TestInsertRecord(t *testing.T) {
	db := connectDB()
	defer db.Close()
	testData := []data{
		{"jeff", "password1", "nil", "false"},
		{"hello@gmail.com", "password2", "nil", "false"},
		{"man", "password3", "hello Pte Ltd", "true"},
	}
	for _, v := range testData {
		err := InsertRecord(db, v.Username, v.Password, v.Company, v.IsCompany)
		if err != nil {
			t.Errorf("Error in inserting information into the database")
		}
	}
}

func TestGetRecord(t *testing.T) {
	db := connectDB()
	defer db.Close()
	_, err := GetRecord(db)
	if err != nil {
		t.Errorf("Error in getting record from database, error is %s", err)
	}
}

func TestSearchRecord(t *testing.T) {
	db := connectDB()
	defer db.Close()
	testData := []data{
		{"jeff", "password1", "nil", "false"},
		{"hello@gmail.com", "password2", "nil", "false"},
		{"man", "password3", "hello Pte Ltd", "true"},
	}
	for _, v := range testData {
		users, err := SearchRecord(db, v.Username)
		if err != nil {
			t.Errorf("Error in searching record, error is %s", err)
		}
		if users.IsCompany != v.IsCompany {
			t.Errorf("Expected company %s, got %s", v.IsCompany, users.IsCompany)
		}
		if users.Company != v.Company {
			t.Errorf("Expected company %s, got %s", v.Company, users.Company)
		}
	}
}
func TestDeleteRecord(t *testing.T) {
	db := connectDB()
	defer db.Close()
	testData := []data{
		{"jeff", "password1", "nil", "false"},
		{"hello@gmail.com", "password2", "nil", "false"},
		{"man", "password3", "hello Pte Ltd", "true"},
	}
	for _, v := range testData {
		err := DeleteRecord(db, v.Username)
		if err != nil {
			t.Errorf("Error in deleting data, error is %s", err)
		}
	}
}
