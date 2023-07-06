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
	ScenarioId                  string `scorm:"pk"`
	Name                        string
	ParentId                    string
	ParentVersionIndex          int
	Level                       int
	CurrentScenarioVersionIndex int
}

type ScenarioVersionTable struct {
	ScenarioVersionId    string `scorm:"pk"`
	ScenarioId           string
	ScenarioVersionIndex int
}

func CreateNewScenario(name string, parent string) *ScenarioTable {
	var newScenario ScenarioTable
	var parentSc *ScenarioTable
	var err error
	var scLevel int = 0
	var parentScvIndex int = 0
	if parent != "" {
		parentSc, err = getScenario(parent)
		if err == nil {
			scLevel = parentSc.Level + 1
			parentScvIndex = parentSc.CurrentScenarioVersionIndex
		}
	}
	if err == nil {
		newScenario = ScenarioTable{ScenarioId: uuid.New().String(), Name: name, ParentId: parent, Level: scLevel, ParentVersionIndex: parentScvIndex, CurrentScenarioVersionIndex: 0}
		newScenarioVersion := ScenarioVersionTable{ScenarioId: newScenario.ScenarioId, ScenarioVersionId: uuid.New().String(), ScenarioVersionIndex: 0}
		Insert(&newScenario, "")
		Insert(&newScenarioVersion, "")
		if parent != "" {
			parentSc.CurrentScenarioVersionIndex++
			newScenarioVersionParent := ScenarioVersionTable{ScenarioId: parent, ScenarioVersionId: uuid.New().String(), ScenarioVersionIndex: parentSc.CurrentScenarioVersionIndex}
			Insert(&newScenarioVersionParent, "")
			Update(parentSc, "")
		}
	}
	return &newScenario
}

func getScenario(scenarioId string) (*ScenarioTable, error) {
	scenario, err := Find(&ScenarioTable{ScenarioId: scenarioId}, "")
	return scenario, err
}

func getScenarioAndCurrentVersion(scenarioId string) (*ScenarioTable, *ScenarioVersionTable, error) {
	scenario, err := Find(&ScenarioTable{ScenarioId: scenarioId}, "")

	if err != nil {
		return scenario, nil, err
	}
	scenarioVersions, err := FindAll(&ScenarioVersionTable{}, "")

	resultIndex := -1
	for i, scv := range scenarioVersions {
		if scv.ScenarioId == scenarioId && scv.ScenarioVersionIndex == scenario.CurrentScenarioVersionIndex {
			resultIndex = i
		}
	}

	return scenario, &scenarioVersions[resultIndex], err
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

func getScenarioAndParentsMap(scenarioId string) (map[string]ScenarioTable, error) {
	result := make(map[string]ScenarioTable)
	scenarios, err := getScenarioAndParents(scenarioId)
	for _, sc := range scenarios {
		result[sc.ScenarioId] = sc
	}
	return result, err
}

func getScenarioVersionsForAllParents(scenarioId string) ([]ScenarioVersionTable, error) {
	var err error
	var result []ScenarioVersionTable
	scenarios, err := FindAll(&ScenarioTable{}, "")
	if err != nil {
		return nil, err
	}
	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].Level > scenarios[j].Level
	})
	scenarioVersions, err := FindAll(&ScenarioVersionTable{}, "")
	sort.Slice(scenarioVersions, func(i, j int) bool {
		return scenarioVersions[i].ScenarioVersionIndex > scenarioVersions[j].ScenarioVersionIndex
	})

	var scvIndex int = 0
	scenarioIdLoop := scenarioId
	for _, sc := range scenarios {
		if sc.ScenarioId == scenarioIdLoop && scenarioId != "" {
			if scenarioIdLoop == scenarioId {
				scvIndex = sc.CurrentScenarioVersionIndex
			}
			for _, scv := range scenarioVersions {
				if scv.ScenarioId == scenarioIdLoop && scv.ScenarioVersionIndex <= scvIndex {
					result = append(result, scv)
				}
			}
			scenarioIdLoop = sc.ParentId
			scvIndex = sc.ParentVersionIndex
		}
	}

	return result, err
}

func createScenarioVersionFieldValuesAndNames(scenarios []ScenarioVersionTable) ([]string, []string) {
	values := []string{}
	fieldNames := []string{}

	for _, sc := range scenarios {
		values = append(values, fmt.Sprintf("'%s'", sc.ScenarioVersionId))
		fieldNames = append(fieldNames, "ScenarioVersionId")
	}
	return values, fieldNames
}
