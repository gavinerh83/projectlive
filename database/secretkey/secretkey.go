package secretkey

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type secret struct {
	Type  string
	Value string
}

//GetKey returns the key
func GetKey(db *sql.DB, secretType string) (string, error) {
	var key string
	var k secret
	err := db.QueryRow("SELECT * FROM secret WHERE Type = ?", secretType).Scan(&k.Type, &k.Value)
	if err != nil {
		return key, err
	}
	key = k.Value
	return key, nil
}
