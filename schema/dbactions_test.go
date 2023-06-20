package schema

import (
	"testing"
)

func TestDbInsert(t *testing.T) {
	insertTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest)
	Insert(&insertTest, "")
}

func TestDbUpdate(t *testing.T) {
	updateTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&updateTest)
	Insert(&updateTest, "")
	updateTest.Nr = 42
	Update(&updateTest, "")
}

func TestDbDelete(t *testing.T) {
	deleteTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&deleteTest)
	Insert(&deleteTest, "")
	Delete(&deleteTest, "")
}

func TestFindAll(t *testing.T) {
	insertTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTest1{Name: "test 2", Flag: false, Nr: 42}
	AutoMigration(&insertTest)
	Insert(&insertTest, "")
	result, _ := FindAll(&insertTest, "")
	if len(result) != 1 {
		t.Fatalf("Expecting 1 row but found %d", len(result))
	}
	Insert(&insertTest2, "")
	result, _ = FindAll(&insertTest, "")
	if len(result) != 2 {
		t.Fatalf("Expecting 2 rows but found %d", len(result))
	}
}

func TestFindAllError(t *testing.T) {
	insertTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	_, err := FindAll(&insertTest, "")
	if err == nil {
		t.Fatalf("Expecting error but not receiving it")
	}
}

func TestDbInsertScenario(t *testing.T) {
	insertTest := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest)
	Insert(&insertTest, "scenario1")

}
