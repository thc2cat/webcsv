package main

import (
	"sort"
)

type ByNameReverse []FileInfo

func (a ByNameReverse) Len() int           { return len(a) }
func (a ByNameReverse) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNameReverse) Less(i, j int) bool { return a[i].Name > a[j].Name } // '>' pour ordre descendant

// SortFilesByNameReverse trie un tableau de FileInfo par le champ Name en ordre inverse ascendant (descendant).
func SortFilesByNameReverse(myfiles []FileInfo) {
	sort.Sort(ByNameReverse(myfiles))
}
