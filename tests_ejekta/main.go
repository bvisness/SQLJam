package main

import (
	"fmt"
	"github.com/bvisness/SQLJam/sqljam/nodes"
)


func main() {
	/*
	inA := sqljam.SqlNodePin{}
	outA := sqljam.SqlNodePin{}
	a := sqljam.SqlNode{}

	a.AddInputPin(&inA)
	a.AddOutputPin(&outA)

	inB := sqljam.SqlNodePin{}
	outB := sqljam.SqlNodePin{}
	b := sqljam.SqlNode{}

	b.AddInputPin(&inB)
	b.AddOutputPin(&outB)

	outA.ConnectTo(&inB)

	fmt.Println(a)
	*/
	src := nodes.SqlNodeTable{TableSource: "SAKILA.FILM"}

	fmt.Println(src)
}