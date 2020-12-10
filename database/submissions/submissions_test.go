package submissions

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var sqluser = "root"
var sqlpassword = "password"

func TestGetDetails(t *testing.T) {
	db := connectDB()
	m, err := GetDetails(db)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(m)
}

func TestGetSpecific(t *testing.T) {
	db := connectDB()
	m, err := GetSpecific(db, "gavinerh@gmail.com")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(m)
}

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
