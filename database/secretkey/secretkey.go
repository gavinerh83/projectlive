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
	stmt, err := db.Prepare("SELECT * FROM secret WHERE Type = ?")
	if err != nil {
		return key, err
	}
	defer stmt.Close()
	results, err := stmt.Query(secretType)
	if err != nil {
		return key, err
	}
	defer results.Close()
	var k secret
	for results.Next() {
		err = results.Scan(&k.Type, &k.Value)
		if err != nil {
			return key, err
		}
	}
	key = k.Value
	return key, nil
}
