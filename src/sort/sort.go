package main

import (
	"fmt"
	"sort"
)

type Person struct {
	age  int
	name string
}

type Group struct {
	persons []Person
	less    func(a, b Person) bool
}

func (g Group) Len() int {
	return len(g.persons)
}

func (g Group) Less(a, b int) bool {
	return g.less(g.persons[a], g.persons[b])
}

func (g Group) Swap(i, j int) {
	g.persons[i], g.persons[j] = g.persons[j], g.persons[i]
}

func main() {
	pSlice := []Person{
		{10, "a"},
		{20, "b"},
		{5, "c"},
		{1, "d"},
	}
	group := Group{pSlice, func(a, b Person) bool { return a.age < b.age }}
	sort.Stable(group)
	fmt.Printf("%v", group.persons)
}
