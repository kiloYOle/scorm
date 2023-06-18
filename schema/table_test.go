package schema

import (
	"testing"
)

func TestTableObjectCreation(t *testing.T) {

	type TableTest struct {
		String01 string `scorm:"pk"`
		Bool01   bool
		Number01 int
		Number02 int8
		Number03 int16
		Number04 int32
		Number05 int64
		Number06 uint
		Number07 uint8
		Number08 uint16
		Number09 uint32
		Number10 uint64
		Number11 float32
		Number12 float64
	}
	table := CreateTableFromStruct(new(TableTest))
	if table.Name != "TableTest" {
		t.Fatalf("Name of the table is not what expected")
	}
	if Count(table.Fields, func(f *Field) bool { return f.Type == String }) != 1 {
		t.Fatalf("Number of string fields is not what expected")
	}
	if Count(table.Fields, func(f *Field) bool { return f.Type == Bool }) != 1 {
		t.Fatalf("Number of bool fields is not what expected")
	}
	if Count(table.Fields, func(f *Field) bool { return f.Type == Int }) != 5 {
		t.Fatalf("Number of int fields is not what expected")
	}
	if Count(table.Fields, func(f *Field) bool { return f.Type == Uint }) != 5 {
		t.Fatalf("Number of uint fields is not what expected")
	}
	if Count(table.Fields, func(f *Field) bool { return f.Type == Float }) != 2 {
		t.Fatalf("Number of float fields is not what expected")
	}
}
