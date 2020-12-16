package binarytree

import "fmt"

type node struct {
	CompanyName string
	Email       string
	left        *node
	right       *node
}

type Binarytree struct {
	Root *node
	size int
}

var s []string

//ReturnSellerInfo contains fields for returning the seller infomation
type ReturnSellerInfo struct {
	CompanyName string
	Email       string
}

func (p *Binarytree) Insert(companyname, email string) error {
	newnode := &node{CompanyName: companyname, Email: email}
	if p.Root == nil {
		p.Root = newnode
	} else { //if not empty
		currentnode := p.Root
		for currentnode != nil {
			if companyname < currentnode.CompanyName {
				//left
				if currentnode.left == nil {
					currentnode.left = newnode
					return nil
				}
				//if left is not empty
				currentnode = currentnode.left
			} else {
				//greater than node, go right
				if currentnode.right == nil {
					currentnode.right = newnode
					return nil
				}
				currentnode = currentnode.right
			}
		}
	}
	return nil
}

func (p *Binarytree) Lookup(companyname string) (ReturnSellerInfo, error) {
	var c ReturnSellerInfo
	if p.Root == nil {
		return c, fmt.Errorf("There are currently no companies signed up")
	}
	currentnode := p.Root
	for currentnode != nil {
		if companyname < currentnode.CompanyName {
			currentnode = currentnode.left
		} else if companyname > currentnode.CompanyName {
			currentnode = currentnode.right
		} else if companyname == currentnode.CompanyName {
			c.CompanyName = currentnode.CompanyName
			c.Email = currentnode.Email
			return c, nil
		}
	}
	return c, fmt.Errorf("Company not found")
}

//ListAllNodes returns the slice of string containing the company names
func (p *Binarytree) ListAllNodes(n *node) []string {
	if n != nil {
		p.ListAllNodes(n.left)
		p.ListAllNodes(n.right)
		s = append(s, n.CompanyName)
	}
	return s
}

//Init create an instance of the binary tree
func Init() *Binarytree {
	tree := &Binarytree{}
	return tree
}

func ResetSlice() {
	if len(s) > 1 {
		s = s[:0]
	}
}
