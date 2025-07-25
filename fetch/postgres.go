package main

import (
	"fmt"
)

type InsertStatement struct {
	Tbl  string
	Cols []string
	Vals [][]any
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
	// ins.addVals(&stmnt)
	return stmnt + ""
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
			// *stmnt += "$"
			if j < (len(r) - 1) {
				*stmnt += ", "
			}
		}
		*stmnt += ")"
		if i < (len(ins.Vals) - 1) {
			*stmnt += ", "
		}
	}
	// *stmnt += ");"
}

func (ins *InsertStatement) addVals(stmnt *string) {
	*stmnt += " values ("
	for i, r := range ins.Vals {
		*stmnt += "("
		for j, v := range r {
			if v == nil {
				v = 0
			}
			*stmnt += fmt.Sprintf("%v", v)
			if j < (len(r) - 1) {
				*stmnt += ", "
			}
		}
		*stmnt += ")"
		if i < (len(ins.Vals) - 1) {
			*stmnt += ", "
		}
	}
	*stmnt += ");"
}
