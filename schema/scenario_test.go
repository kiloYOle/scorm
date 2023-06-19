package schema

import (
	"testing"
)

func TestCreateScenario(t *testing.T) {
	CreateScenarioTables()
	CreateNewScenario("Test scenario", "")
}

func TestCreateChildScenario(t *testing.T) {
	CreateScenarioTables()
	scenario := CreateNewScenario("Test scenario", "")
	CreateNewScenario("Test scenario 2", scenario.ScenarioId)
}
