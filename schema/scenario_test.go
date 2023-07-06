package schema

import (
	"testing"
)

func TestCreateScenario(t *testing.T) {
	CreateScenarioTables()
	CreateNewScenario("Test scenario", "")
	result, _ := FindAll(&ScenarioTable{}, "")
	if len(result) != 1 {
		t.Fatalf("Expecting 1 scenario but found %d", len(result))
	}
}

func TestCreateChildScenarioFail(t *testing.T) {
	result, _ := FindAll(&ScenarioTable{}, "")
	scenariosStart := len(result)
	CreateScenarioTables()
	CreateNewScenario("Test scenario", "not existing")
	result, _ = FindAll(&ScenarioTable{}, "")
	if len(result)-scenariosStart != 0 {
		t.Fatalf("Expecting 0 scenarios but found %d", len(result)-scenariosStart)
	}
}

func TestCreateChildScenario(t *testing.T) {
	result, _ := FindAll(&ScenarioTable{}, "")
	result_scv, _ := FindAll(&ScenarioVersionTable{}, "")
	scenariosStart := len(result)
	scenarioVsStart := len(result_scv)
	CreateScenarioTables()
	scenario := CreateNewScenario("Test scenario", "")
	CreateNewScenario("Test scenario 2", scenario.ScenarioId)
	result, _ = FindAll(&ScenarioTable{}, "")
	result_scv, _ = FindAll(&ScenarioVersionTable{}, "")
	if len(result)-scenariosStart != 2 {
		t.Fatalf("Expecting 2 scenarios but found %d", len(result)-scenariosStart)
	}
	if len(result_scv)-scenarioVsStart != 3 {
		t.Fatalf("Expecting 3 scenarios versions but found %d", len(result_scv)-scenarioVsStart)
	}
}

func TestGetScenarioParents(t *testing.T) {
	CreateScenarioTables()
	scenario := CreateNewScenario("Test scenario", "")
	scenario2 := CreateNewScenario("Test scenario 2", scenario.ScenarioId)
	scenario3 := CreateNewScenario("Test scenario 3", scenario2.ScenarioId)
	result, _ := getScenarioAndParents(scenario.ScenarioId)
	if len(result) != 1 {
		t.Fatalf("Expecting 1 scenario but found %d", len(result))
	}
	result, _ = getScenarioAndParents(scenario2.ScenarioId)
	if len(result) != 2 {
		t.Fatalf("Expecting 2 scenarios but found %d", len(result))
	}
	result, _ = getScenarioAndParents(scenario3.ScenarioId)
	if len(result) != 3 {
		t.Fatalf("Expecting 3 scenarios but found %d", len(result))
	}
}

func TestGetScenarioAndVersion(t *testing.T) {
	CreateScenarioTables()
	scenario := CreateNewScenario("Test scenario", "")
	scenario2 := CreateNewScenario("Test scenario 2", scenario.ScenarioId)
	result_sc, result_scv, _ := getScenarioAndCurrentVersion(scenario.ScenarioId)
	if result_sc != nil && result_scv != nil {
		t.Fatalf("Expecting 1 of each record but not found all")
	}
	result_sc, result_scv, _ = getScenarioAndCurrentVersion(scenario2.ScenarioId)
	if result_sc != nil && result_scv != nil {
		t.Fatalf("Expecting 1 of each record but not found all")
	}
}
