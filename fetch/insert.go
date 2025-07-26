package main

import (
	"fmt"
)

type InsertStatement struct {
	Tbl     string
	PrimKey string // define like "key" or "key1, key2"
	Cols    []string
	Vals    [][]any
}

// flatten [][]any to []any
func (ins *InsertStatement) FlattenVals() []any {
	var valsFlat []any
	for _, r := range ins.Vals {
		valsFlat = append(valsFlat, r...)
	}
	return valsFlat
}

func (ins *InsertStatement) Build() string {
	stmnt := fmt.Sprintf("insert into %s (", ins.Tbl)
	ins.addCols(&stmnt)
	ins.addValsPlHldr(&stmnt)
	return fmt.Sprintf("%s on conflict (%s) do nothing", stmnt, ins.PrimKey)
}

func (ins *InsertStatement) addCols(stmnt *string) {
	for i, c := range ins.Cols {
		*stmnt += c
		if i < (len(ins.Cols) - 1) {
			*stmnt += ", "
		}
	}
	*stmnt += ")"
}

func (ins *InsertStatement) addValsPlHldr(stmnt *string) {
	*stmnt += " values "
	for i, r := range ins.Vals {
		*stmnt += "("
		for j := range r {
			*stmnt += fmt.Sprintf("$%d", i*len(r)+(j+1))
			if j < (len(r) - 1) {
				*stmnt += ", "
			}
		}
		*stmnt += ")"
		if i < (len(ins.Vals) - 1) {
			*stmnt += ", "
		}
	}
}
