package schema

import (
	"testing"
)

type TableTest1 struct {
	String01_1 string `scorm:"pk"`
	Bool01_1   bool
	Number01_1 int
	Number02_1 int8 `scorm:"pk"`
	Number03_1 int16
	Number04_1 int32
	Number05_1 int64
	Number06_1 uint
	Number07_1 uint8
	Number08_1 uint16
	Number09_1 uint32
	Number10_1 uint64
	Number11_1 float32
	Number12_1 float64
}

type TableTest2 struct {
	String01_2 string
	Bool01_2   bool
	Number01_2 int
	Number02_2 int8 `scorm:"pk"`
	Number03_2 int16
	Number04_2 int32
	Number05_2 int64 `scorm:"pk"`
	Number06_2 uint
	Number07_2 uint8
	Number08_2 uint16 `scorm:"pk"`
	Number09_2 uint32
	Number10_2 uint64
	Number11_2 float32
	Number12_2 float64
}

type ScenarioTableTest1 struct {
	Scenario
	String01_2 string `scorm:"pk"`
	Bool01_2   bool
	Number01_2 int
}

func TestAutomigration(t *testing.T) {
	OpenTestDB()
	AutoMigration(&TableTest1{}, &TableTest2{})
	AutoMigration()
}

func TestScenarioAutomigration(t *testing.T) {
	OpenTestDB()
	AutoMigration(&ScenarioTableTest1{})
}
