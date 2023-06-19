package schema

import (
	"github.com/google/uuid"
)

type Scenario struct {
	ScenarioVersionId string `scorm:"pk"`
	IsDeleted         bool
}

type ScenarioTable struct {
	ScenarioId string `scorm:"pk"`
	Name       string
	ParentId   string
	Level      int
}

type ScenarioVersionTable struct {
	ScenarioVersionId       string `scorm:"pk"`
	ScenarioId              string
	PreviousScenarioVersion string
}

func CreateNewScenario(name string, parent string) *ScenarioTable {
	newScenario := ScenarioTable{ScenarioId: uuid.New().String(), Name: name, ParentId: parent, Level: 0}
	newScenarioVersion := ScenarioVersionTable{ScenarioId: newScenario.ScenarioId, ScenarioVersionId: uuid.New().String(), PreviousScenarioVersion: ""}
	Insert(&newScenario, "")
	Insert(&newScenarioVersion, "")
	return &newScenario
}
