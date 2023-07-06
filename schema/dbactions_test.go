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

func TestFind(t *testing.T) {
	insertTest := InsertTest1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTest1{Name: "test 2", Flag: false, Nr: 42}
	AutoMigration(&insertTest)
	Insert(&insertTest, "")
	result, _ := Find(&insertTest, "")
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
	Insert(&insertTest2, "")
	result, _ = Find(&insertTest, "")
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
}

func TestDbInsertScenario(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTestScenario1{Name: "test 2", Flag: true, Nr: 42}
	insertTest3 := InsertTestScenario1{Name: "test 3", Flag: true, Nr: 13}
	AutoMigration(&insertTest1)
	scenario1 := CreateNewScenario("scenario1", "")
	Insert(&insertTest1, scenario1.ScenarioId)
	result, _ := FindAll(&insertTest1, scenario1.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 row but found %d", len(result))
	}
	scenario2 := CreateNewScenario("scenario2", scenario1.ScenarioId)
	Insert(&insertTest2, scenario2.ScenarioId)
	Insert(&insertTest3, scenario1.ScenarioId)
	result, _ = FindAll(&insertTest1, scenario1.ScenarioId)
	if len(result) != 2 {
		t.Fatalf("Expecting 2 rows but found %d", len(result))
	}
	result, _ = FindAll(&insertTest1, scenario2.ScenarioId)
	if len(result) != 2 {
		t.Fatalf("Expecting 2 rows but found %d", len(result))
	}
}

func TestDbFindAllScenario(t *testing.T) {
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

func TestDbFindScenario(t *testing.T) {
	insertTest := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest)
	Insert(&insertTest, "not existing")
	result, _ := Find(&insertTest, "")
	if result != nil {
		t.Fatalf("Expecting no row but found one")
	}
	scenario := CreateNewScenario("scenario1", "")
	Insert(&insertTest, scenario.ScenarioId)
	result, _ = Find(&insertTest, scenario.ScenarioId)
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
}

func TestDbFindAllScenarioParent(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTestScenario1{Name: "test 2", Flag: true, Nr: 13}
	insertTest3 := InsertTestScenario1{Name: "test 3", Flag: true, Nr: 14}
	AutoMigration(&insertTest1)
	scenario_parent := CreateNewScenario("scenario1", "")
	scenario_child := CreateNewScenario("scenario2", scenario_parent.ScenarioId)
	scenario_child2 := CreateNewScenario("scenario3", scenario_child.ScenarioId)
	Insert(&insertTest1, scenario_parent.ScenarioId)
	Insert(&insertTest2, scenario_child.ScenarioId)
	Insert(&insertTest3, scenario_child2.ScenarioId)
	result, _ := FindAll(&insertTest1, scenario_parent.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 row but found %d", len(result))
	}
	result, _ = FindAll(&insertTest1, scenario_child.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 rows but found %d", len(result))
	}
	result, _ = FindAll(&insertTest1, scenario_child2.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 rows but found %d", len(result))
	}
}

func TestDbFindScenarioParent(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	insertTest2 := InsertTestScenario1{Name: "test 2", Flag: true, Nr: 13}
	insertTest3 := InsertTestScenario1{Name: "test 3", Flag: true, Nr: 14}
	AutoMigration(&insertTest1)
	scenario_parent := CreateNewScenario("scenario1", "")
	Insert(&insertTest1, scenario_parent.ScenarioId)
	scenario_child := CreateNewScenario("scenario2", scenario_parent.ScenarioId)
	scenario_child2 := CreateNewScenario("scenario3", scenario_child.ScenarioId)
	Insert(&insertTest2, scenario_child.ScenarioId)
	Insert(&insertTest3, scenario_child2.ScenarioId)
	result, _ := Find(&insertTest1, scenario_parent.ScenarioId)
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
	result, _ = Find(&insertTest1, scenario_child.ScenarioId)
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
	result, _ = Find(&insertTest1, scenario_child2.ScenarioId)
	if result == nil {
		t.Fatalf("Expecting 1 row but found none")
	}
	result, _ = Find(&insertTest3, scenario_child.ScenarioId)
	if result != nil {
		t.Fatalf("Expecting no row but found one")
	}
}

func TestDbUpdateScenario(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest1)
	scenario_parent := CreateNewScenario("scenario1", "")
	Insert(&insertTest1, scenario_parent.ScenarioId)
	scenario_child := CreateNewScenario("scenario2", scenario_parent.ScenarioId)
	scenario_child2 := CreateNewScenario("scenario3", scenario_child.ScenarioId)
	insertTest1.Nr = 21
	Update(&insertTest1, scenario_parent.ScenarioId)
	insertTest1.Nr = 22
	Update(&insertTest1, scenario_child.ScenarioId)
	insertTest1.Nr = 23
	Update(&insertTest1, scenario_child2.ScenarioId)
	result, _ := Find(&insertTest1, scenario_parent.ScenarioId)
	if result.Nr != 21 {
		t.Fatalf("Expecting 21 but found %d", result.Nr)
	}
	result, _ = Find(&insertTest1, scenario_child.ScenarioId)
	if result.Nr != 22 {
		t.Fatalf("Expecting 22 row but found %d", result.Nr)
	}
	result, _ = Find(&insertTest1, scenario_child2.ScenarioId)
	if result.Nr != 23 {
		t.Fatalf("Expecting 23 row but found %d", result.Nr)
	}
}

func TestDbUpdateScenarioKeepParentOriginal(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest1)
	scenario_parent := CreateNewScenario("scenario1", "")
	scenario_child := CreateNewScenario("scenario2", scenario_parent.ScenarioId)
	Insert(&insertTest1, scenario_parent.ScenarioId)
	insertTest1.Nr = 21
	Update(&insertTest1, scenario_child.ScenarioId)
	result, _ := Find(&insertTest1, scenario_parent.ScenarioId)
	if result.Nr != 12 {
		t.Fatalf("Expecting 12 but found %d", result.Nr)
	}
	result, _ = Find(&insertTest1, scenario_child.ScenarioId)
	if result.Nr != 21 {
		t.Fatalf("Expecting 21 row but found %d", result.Nr)
	}
}

func TestDbUpdateScenarioUpdateParentKeepChildOriginal(t *testing.T) {
	insertTest1 := InsertTestScenario1{Name: "test", Flag: true, Nr: 12}
	AutoMigration(&insertTest1)
	scenario_parent := CreateNewScenario("scenario1", "")
	Insert(&insertTest1, scenario_parent.ScenarioId)
	scenario_child := CreateNewScenario("scenario2", scenario_parent.ScenarioId)
	insertTest1.Nr = 21
	Update(&insertTest1, scenario_parent.ScenarioId)
	result, _ := Find(&insertTest1, scenario_parent.ScenarioId)
	if result.Nr != 21 {
		t.Fatalf("Expecting 21 but found %d", result.Nr)
	}
	result, _ = Find(&insertTest1, scenario_child.ScenarioId)
	if result.Nr != 12 {
		t.Fatalf("Expecting 12 but found %d", result.Nr)
	}
}
