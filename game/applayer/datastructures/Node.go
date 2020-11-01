package datastructures

import (
	"fmt"
)

type Node struct {
	Value interface{}
}

func (n *Node) String() string {
	return fmt.Sprint(n.Value)
}
