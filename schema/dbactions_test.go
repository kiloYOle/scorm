package schema

import (
	"testing"
)

func TestDbInsert(t *testing.T) {
	insertTest := InsertTest1{name: "test", flag: true, nr: 12}
	AutoMigration(&insertTest)
	Insert(&insertTest)
}

func TestDbUpdate(t *testing.T) {
	updateTest := InsertTest1{name: "test", flag: true, nr: 12}
	AutoMigration(&updateTest)
	Insert(&updateTest)
	updateTest.nr = 42
	Update(&updateTest)
}

func TestDbDelete(t *testing.T) {
	deleteTest := InsertTest1{name: "test", flag: true, nr: 12}
	AutoMigration(&deleteTest)
	Insert(&deleteTest)
	Delete(&deleteTest)
}
