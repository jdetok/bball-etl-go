package main

import (
	"fmt"
)

type InsertStatement struct {
	Tbl     string
	PrimKey string // define like "key" or "key1, key2"
	Cols    []string
	Vals    []any
	ValSets [][]any
}

// flatten [][]any to []any
func (ins *InsertStatement) FlattenVals() []any {
	var valsFlat []any
	for _, r := range ins.ValSets {
		valsFlat = append(valsFlat, r...)
	}
	return valsFlat
}

func Flatten(set [][]any) []any {
	var valsFlat []any
	for _, r := range set {
		valsFlat = append(valsFlat, r...)
	}
	return valsFlat
}

func (ins *InsertStatement) ChunkVals() {
	const PG_MAX int = 65000
	var totSets int = len(ins.ValSets)
	var totVals int = len(ins.FlattenVals())
	var valsPer int = len(ins.ValSets[0])
	var setsFit int = PG_MAX / valsPer
	var numChunks int = (totVals + PG_MAX - 1) / PG_MAX

	chunkIdx := make([][2]int, 0, numChunks)
	for i := range numChunks {
		start := i * setsFit
		end := start + setsFit
		if end > totSets {
			fmt.Printf("end (%d) > total sets (%d)\n", end, totSets)
			end = totSets
		}
		chunkIdx = append(chunkIdx, [2]int{start, end})
	}

	var chunks [][][]any
	for i, c := range chunkIdx {
		var valChunk [][]any = ins.ValSets[c[0]:c[1]]
		fmt.Printf("chunk %d - sets: %d | vals: %d\n", i, len(valChunk),
			len(Flatten(valChunk)))
		chunks = append(chunks, valChunk)
	}
	fmt.Println(len(chunks))
}

/*
postgres will only insert 65535 vals at a time. count number of vals & split up
into chunks that can be inserted
*/
func (ins *InsertStatement) ChunkValsTest() {
	const PG_MAX_VALS int = 65000
	var totValSets int = len(ins.ValSets)
	var totVals int = len(ins.FlattenVals())
	var valsPer int = len(ins.ValSets[0])
	var numSetsFit int = PG_MAX_VALS / valsPer
	var numValsFit int = numSetsFit * valsPer
	var numSetsRem int = totValSets - numSetsFit
	var numValsRem int = numSetsRem * valsPer
	var numChunks int = (totVals + PG_MAX_VALS - 1) / PG_MAX_VALS
	fmt.Println("numchunks:", numChunks)

	fmt.Println("num of valsets:", totValSets)
	fmt.Println("total vals:", totVals)
	fmt.Println("vals in one valset:", valsPer)
	fmt.Printf("num rowsets with vals < max (%d): %d:\n", PG_MAX_VALS, numSetsFit)
	fmt.Println("num vals in rowsets that fit under max:", numValsFit)
	fmt.Println("num rowsets remaining:", numSetsRem)
	fmt.Println("num vals in remaining rowsets:", numValsRem)

	var chunks [][]int
	for i := range numChunks {
		var start int = i * numSetsFit
		var end int
		if i == 0 {
			end = start + numSetsFit
		} else {
			end = totValSets
		}
		var chunk = []int{start, end}
		chunks = append(chunks, chunk)
	}
	fmt.Println(chunks)

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
	for i, r := range ins.ValSets {
		*stmnt += "("
		for j := range r {
			// postgres uses 1-type placeholders, i*rows + idx of current val
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
