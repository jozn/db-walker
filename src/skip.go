package src

import "strings"

var notMakeTableType = []string{"user"}

func skipTableModel(table string) bool {
	t := strings.ToLower(table)
	for _, ent := range notMakeTableType {
		if ent == t {
			return true
		}
	}
	return false
}
