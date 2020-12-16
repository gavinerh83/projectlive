package hashtable

import (
	"testing"
)

//Data consist of fields for testing the hashtable
type Data struct {
	key   string
	value string
}

type transaction struct {
	key   string
	value string
	temp  string
}

var hashTable = Init()

func TestInsert(t *testing.T) {
	testData := []Data{
		{"Randy", "best"},
		{"Tinkerbell", "Peter"},
		{"Tom", "Jerry"},
		{"Randy", "Foam"},
	}

	for _, v := range testData {
		err := hashTable.Insert(v.key, v.value)
		if err != nil {
			if err.Error() == "User already exist" {
				t.Logf("Test for adding existing user passed")
			} else {
				t.Errorf("Expected nil got %s", err.Error())
			}
		}
	}
}

func TestSearch(t *testing.T) {
	testData := []Data{
		{"Randy", "best"},
		{"Tinkerbell", "Peter"},
		{"Tom", "Jerry"},
		{"Captain", "Troy"},
	}
	for _, v := range testData {
		s, err := hashTable.Search(v.key)
		if err != nil {
			if err.Error() != "User not found" {
				t.Errorf("Expected %s got %s", "User not found", err.Error())
			}
		} else {
			if s != v.value {
				t.Errorf("Expected %s got %s", v.value, s)
			}
		}
	}
}

func TestInsertTransaction(t *testing.T) {
	testData := []transaction{
		{"Tinkerbell", "Peter", "hello"},
		{"Derrick", "Johnny", "Lee"},
		{"Henrick", "Que", "Loft"},
		{"Welden", "Robot", "Auto"},
	}
	for _, v := range testData {
		err := hashTable.InsertTransaction(v.key, v.value, v.temp)
		if err != nil {
			if err.Error() == "User already exist" {
				t.Logf("Testing for presence of user passed")
			} else {
				t.Errorf(err.Error())
			}
		}
	}
}

func TestSearchTransaction(t *testing.T) {
	testData := []transaction{
		{"Tiny", "Peter", ""},
		{"Derrick", "Johnny", "Lee"},
		{"Henrick", "Que", "Loft"},
		{"Welden", "Robot", "Auto"},
	}
	for _, v := range testData {
		s, err := hashTable.SearchTransaction(v.key)
		if err != nil {
			if err.Error() == "User not found" {
				t.Logf("Testing for presence of user passed")
			} else {
				t.Errorf(err.Error())
			}
		}
		if s != v.temp {
			t.Errorf("Expected %s got %s", v.temp, s)
		}
	}
}
func TestDelete(t *testing.T) {
	testData := []Data{
		{"Randy", "best"},
		{"Tinkerbell", "Peter"},
		{"Tom", "Jerry"},
		{"Captain", "Troy"},
	}
	for _, v := range testData {
		_, err := hashTable.Search(v.key)
		isDeleted := hashTable.Delete(v.key)
		if err != nil && isDeleted == false {
			t.Logf("Test for deleting non-existent user passed")
		} else {
			if isDeleted != true {
				t.Errorf("Expected %t got %t", true, isDeleted)
			}
		}
	}
}
