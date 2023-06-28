package schema

import (
	"fmt"
	"sort"

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
	var newScenario ScenarioTable
	var parentSc *ScenarioTable
	var err error
	var scLevel int = 0
	if parent != "" {
		parentSc, err = getScenario(parent)
		if err == nil {
			scLevel = parentSc.Level + 1
		}
	}
	if err == nil {
		newScenario = ScenarioTable{ScenarioId: uuid.New().String(), Name: name, ParentId: parent, Level: scLevel}
		newScenarioVersion := ScenarioVersionTable{ScenarioId: newScenario.ScenarioId, ScenarioVersionId: uuid.New().String(), PreviousScenarioVersion: ""}
		Insert(&newScenario, "")
		Insert(&newScenarioVersion, "")
	}
	return &newScenario
}

func getScenario(scenarioId string) (*ScenarioTable, error) {
	var scenario *ScenarioTable
	var err error
	scenarios, err := Find(&ScenarioTable{ScenarioId: scenarioId}, "")
	if len(scenarios) > 0 {
		scenario = &scenarios[0]
	} else {
		err = fmt.Errorf("cannot find scenario with id %s", scenarioId)
	}
	return scenario, err
}

func getScenarioAndParents(scenarioId string) ([]ScenarioTable, error) {
	var err error
	var result []ScenarioTable
	scenarios, err := FindAll(&ScenarioTable{ScenarioId: scenarioId}, "")
	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].Level > scenarios[j].Level
	})
	for _, sc := range scenarios {
		if sc.ScenarioId == scenarioId && scenarioId != "" {
			result = append(result, sc)
			scenarioId = sc.ParentId
		}
	}

	return result, err
}

func createScenarioFieldValuesAndNames(scenarios []ScenarioTable) ([]string, []string) {
	values := []string{}
	fieldNames := []string{}

	for _, sc := range scenarios {
		values = append(values, fmt.Sprintf("'%s'", sc.ScenarioId))
		fieldNames = append(fieldNames, "scenarioId")
	}
	return values, fieldNames
}
