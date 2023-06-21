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

func TestCreateChildScenario(t *testing.T) {
	result, _ := FindAll(&ScenarioTable{}, "")
	scenariosStart := len(result)
	CreateScenarioTables()
	scenario := CreateNewScenario("Test scenario", "")
	CreateNewScenario("Test scenario 2", scenario.ScenarioId)
	result, _ = FindAll(&ScenarioTable{}, "")
	if len(result)-scenariosStart != 2 {
		t.Fatalf("Expecting 2 scenarios but found %d", len(result)-scenariosStart)
	}
}
