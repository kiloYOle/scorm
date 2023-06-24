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
	// will fail as db table not created
	_, err := FindAll(&insertTest, "")
	if err == nil {
		t.Fatalf("Expecting an error but not receiving it")
	}
}

func TestDbInsertScenario(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTestScenario1{Name: "test 2", Flag: true, Nr: 42}
	AutoMigration(&insertTest1)
	scenario1 := CreateNewScenario("scenario1", "")
	scenario2 := CreateNewScenario("scenario2", "")
	Insert(&insertTest1, scenario1.ScenarioId)
	Insert(&insertTest1, scenario2.ScenarioId)
	Insert(&insertTest2, scenario2.ScenarioId)
	result, _ := FindAll(&insertTest1, scenario1.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 row but found %d", len(result))
	}
	result, _ = FindAll(&insertTest1, scenario2.ScenarioId)
	if len(result) != 2 {
		t.Fatalf("Expecting 2 rows but found %d", len(result))
	}
}

func TestDbFindScenario(t *testing.T) {
	insertTest := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest)
	Insert(&insertTest, "not existing")
	result, _ := FindAll(&insertTest, "")
	if len(result) != 0 {
		t.Fatalf("Expecting 0 rows but found %d", len(result))
	}
	scenario := CreateNewScenario("scenario1", "")
	Insert(&insertTest, scenario.ScenarioId)
	result, _ = FindAll(&insertTest, scenario.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 row but found %d", len(result))
	}
}
