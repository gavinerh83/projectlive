//Package hashtable contains the structure of implementing a hash table
package hashtable

import (
	"fmt"
)

const arraySize = 10

//HashTable contains the structure of the hashTable, which is an array
type HashTable struct {
	array [arraySize]*Bucket
}

//Bucket contains the structure of the individual location of the addresses
type Bucket struct {
	head *Node
	size int
}

//Node contain the fields that an instance of the hash table should have
type Node struct {
	Username    string
	Value       interface{}
	Password    string
	Company     string
	IsCompany   string
	TempStorage string
	next        *Node
}

//hash generates a hash from the inserted key
func hash(key string) int {
	sum := 0
	for _, v := range key {
		sum += int(v)
	}
	return sum % arraySize
}

//Insert inserts into the hashtable
func (h *HashTable) Insert(key string, value interface{}) error {
	index := hash(key)
	return h.array[index].insert(key, value)
}

//Search looks for the key in the hashtable
func (h *HashTable) Search(key string) (interface{}, error) {
	index := hash(key)
	return h.array[index].search(key)
}

// Delete removes the item from the hash function which will call the delete from the linkedlist
func (h *HashTable) Delete(key string) bool {
	index := hash(key)
	return h.array[index].delete(key)
}

//insert inserts the key in the linked list
func (p *Bucket) insert(k string, v interface{}) error {
	//first check if the key already exist
	if p.searchPresence(k) == true {
		//if item exist
		return fmt.Errorf("User already exist")
	}
	//inserted item becomes the new head
	newnode := &Node{Username: k, Value: v}
	newnode.next = p.head
	p.head = newnode
	return nil
}

//ssearchPresence search for key and return if key is found
func (p *Bucket) searchPresence(k string) bool {
	currentnode := p.head
	//keep looping the list until you find the item
	for currentnode != nil {
		if currentnode.Username == k {
			return true
		} else {
			currentnode = currentnode.next
		}
	}
	return false
}

//search finds the key from the linkedlist from determined array index and returns the value tagged to key
func (p *Bucket) search(k string) (interface{}, error) {
	currentnode := p.head
	//keep looping the list until you find the item
	for currentnode != nil {
		if currentnode.Username == k {
			return currentnode.Value, nil
		} else {
			currentnode = currentnode.next
		}
	}
	return "", fmt.Errorf("User not found")
}

func (p *Bucket) delete(k string) bool {
	if p.head == nil {
		return false
	}
	//check if the head is the key
	if p.head.Username == k {
		if p.head.next == nil {
			p.head = nil
			return true
		}
		p.head = p.head.next
		return true
	}
	currentnode := p.head
	for currentnode != nil {
		if currentnode.next.Username == k {
			currentnode.next = currentnode.next.next
			return true
		} else {
			currentnode = currentnode.next
		}
	}
	return false
}

//Init initialise the hashtable and create a bucket linkedlist in each slot of memory
func Init() *HashTable {
	result := &HashTable{}
	for i := range result.array {
		result.array[i] = &Bucket{}
	}
	return result
}
