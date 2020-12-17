package binarytree

import (
	"testing"
)

type data struct {
	CompanyName string
	Email       string
}

var testData = []data{
	{CompanyName: "company3", Email: "test2@gmail.com"},
	{CompanyName: "company1", Email: "test@gmail.com"},
	{CompanyName: "company2", Email: "test3@gmail.com"},
	{CompanyName: "company5", Email: "test4@gmail.com"},
	{CompanyName: "company7", Email: "test5@gmail.com"},
	{CompanyName: "company9", Email: "test6@gmail.com"},
}

var tree = Init()

func TestInsert(t *testing.T) {
	for _, v := range testData {
		err := tree.Insert(v.CompanyName, v.Email)
		if err != nil {
			t.Errorf("Expected error nil got %s", err.Error())
		}
	}
}

func TestLookup(t *testing.T) {
	for _, v := range testData {
		sinfo, err := tree.Lookup(v.CompanyName)
		if err != nil {
			t.Errorf("Expected error nil, got %s", err.Error())
		}
		if sinfo.Email != v.Email && sinfo.CompanyName != v.CompanyName {
			t.Errorf("Expected %s and %s got %s %s", v.Email, v.CompanyName, sinfo.Email, sinfo.CompanyName)
		}
	}
}

func TestListAllNodes(t *testing.T) {
	s := tree.ListAllNodes(tree.Root)
	if len(s) != 6 {
		t.Errorf("Expected 6 items, only got %d", len(s))
	}
}
