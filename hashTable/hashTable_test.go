package hashtable

import (
	"testing"
)

//Data consist of fields for testing the hashtable
type Data struct {
	key   string
	value string
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
			if err.Error() == "User not found" {
				t.Logf("Test for searching non-existent user passed")
			}
		} else {
			if s != v.value {
				t.Errorf("Expected %s got %s", v.value, s)
			}
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
