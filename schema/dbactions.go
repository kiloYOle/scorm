package schema

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func Insert(value interface{}, scenarioId string) error {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	if table.IsScenarioBased {
		_, scenarioVersion, err := getScenarioAndCurrentVersion(scenarioId)
		if err == nil {
			addScenarioFields(&values, &fNames, scenarioVersion.ScenarioVersionId)
		} else {
			fmt.Println(err)
			return nil
		}
	}
	insertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(fNames, ","), strings.Join(values, ","))
	fmt.Println(insertString)
	_, err := DB.Exec(insertString)
	if err != nil {
		fmt.Println("Insert error", err)
	}
	return nil
}

func Update[T any](value *T, scenarioId string) {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	pkValues, pkNames := CreatePKFieldValuesAndNames(table, value)
	insertScenarioRow := false
	var setFieldString []string
	for i, fName := range fNames {
		setFieldString = append(setFieldString, fName+"="+values[i])
	}
	whereClause := createWhereClauseAnd(&pkValues, &pkNames)
	if table.IsScenarioBased {
		scenario, _ := getScenario(scenarioId)
		//scenarioMap, _ := getScenarioAndParentsMap(scenario)
		scenarioVersions, _ := getScenarioVersionsForAllParents(scenarioId)
		_, resultScenarioVersion, err := findFromAllScenarios(value, scenarioId)
		if err != nil {
			return
		}
		scvIndex := -1
		for i, scv := range resultScenarioVersion {
			if scv.ScenarioId == scenarioId && scv.ScenarioVersionIndex == scenario.CurrentScenarioVersionIndex {
				scvIndex = i
				break
			}
		}

		if scvIndex < 0 {
			insertScenarioRow = true
			var scenarioVersionId string
			for _, scv := range scenarioVersions {
				if scv.ScenarioId == scenarioId && scv.ScenarioVersionIndex == scenario.CurrentScenarioVersionIndex {
					scenarioVersionId = scv.ScenarioVersionId
				}
			}
			addScenarioFields(&values, &fNames, scenarioVersionId)
		} else {
			scValues, scNames := createScenarioVersionFieldValuesAndNames(resultScenarioVersion[scvIndex:1])
			whereClause_sc := createWhereClauseOr(&scValues, &scNames)
			whereClause = fmt.Sprintf("(%s) AND (%s)", whereClause, whereClause_sc)
		}
	}
	var execString string
	if insertScenarioRow {
		execString = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(fNames, ","), strings.Join(values, ","))
	} else {
		execString = fmt.Sprintf("UPDATE %s SET %s WHERE %s", table.Name, strings.Join(setFieldString, ","), whereClause)
	}
	fmt.Println(execString)
	_, err := DB.Exec(execString)
	if err != nil {
		fmt.Println(err)
	}
}

func Delete(value interface{}, scenario string) {
	table := CreateTableFromStruct(value)
	pkValues, pkNames := CreatePKFieldValuesAndNames(table, value)
	var whereString []string
	for i, fName := range pkNames {
		whereString = append(whereString, fName+"="+pkValues[i])
	}
	deleteString := fmt.Sprintf("DELETE FROM %s WHERE %s", table.Name, strings.Join(whereString, " AND "))
	fmt.Println(deleteString)
	_, err := DB.Exec(deleteString)
	if err != nil {
		fmt.Println(err)
	}
}

func Find[T any](value *T, scenario string) (*T, error) {
	result, resultScenarioVersion, err := findFromAllScenarios(value, scenario)
	if err != nil {
		return nil, err
	}
	scenarioMap, err := getScenarioAndParentsMap(scenario)
	if len(result) < 1 {
		fmt.Println("no record found in Find")
		return nil, fmt.Errorf("no record found")
	}
	if scenario == "" {
		return &result[0], err
	}
	fmt.Println("find result")
	fmt.Println(result)
	var maxLevel, maxIndex, index int = -1, -1, -1
	for i := range result {
		if scenarioMap[resultScenarioVersion[i].ScenarioId].Level > maxLevel {
			maxLevel = scenarioMap[resultScenarioVersion[i].ScenarioId].Level
			maxIndex = resultScenarioVersion[i].ScenarioVersionIndex
			index = i
		}
		if scenarioMap[resultScenarioVersion[i].ScenarioId].Level == maxLevel && resultScenarioVersion[i].ScenarioVersionIndex > maxIndex {
			maxIndex = resultScenarioVersion[i].ScenarioVersionIndex
			index = i
		}
	}

	return &result[index], err
}

func findFromAllScenarios[T any](value *T, scenario string) ([]T, []ScenarioVersionTable, error) {
	var result []T
	var resultScenarioVersions []ScenarioVersionTable
	var scenarioVersions []ScenarioVersionTable
	table := CreateTableFromStruct(value)
	values, names := CreatePKFieldValuesAndNames(table, value)
	whereClause := createWhereClauseAnd(&values, &names)
	if table.IsScenarioBased {
		scenarioVersions, _ = getScenarioVersionsForAllParents(scenario)
		scvValues, scvNames := createScenarioVersionFieldValuesAndNames(scenarioVersions)
		whereClause_sc := createWhereClauseOr(&scvValues, &scvNames)
		whereClause = fmt.Sprintf("(%s) AND (%s)", whereClause, whereClause_sc)
	}
	fmt.Printf("SELECT * FROM %s WHERE %s\n", table.NameDB, whereClause)
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s", table.NameDB, whereClause))
	if err != nil {
		fmt.Println(err)
		return result, nil, err
	}
	defer rows.Close()

	result, scenarioVersionIds, err := convertRowsToObject(value, rows)
	fmt.Println(scenarioVersionIds)
	if table.IsScenarioBased {
		for _, s := range scenarioVersionIds {
			for _, scv := range scenarioVersions {
				if scv.ScenarioVersionId == s {
					resultScenarioVersions = append(resultScenarioVersions, scv)
				}
			}

		}
	}

	return result, resultScenarioVersions, err
}

func FindAll[T any](value *T, scenario string) ([]T, error) {
	table := CreateTableFromStruct(value)
	var whereClause string
	if table.IsScenarioBased {
		scenarioVersions, _ := getScenarioVersionsForAllParents(scenario)
		fmt.Printf("find all nr scv %d\n", len(scenarioVersions))
		scValues, scNames := createScenarioVersionFieldValuesAndNames(scenarioVersions)
		whereClause = createWhereClauseOr(&scValues, &scNames)
	}
	var rows *sql.Rows
	var err error
	if whereClause == "" {
		fmt.Printf("SELECT * FROM %s\n", table.NameDB)
		rows, err = DB.Query(fmt.Sprintf("SELECT * FROM %s", table.NameDB))
	} else {
		fmt.Printf("SELECT * FROM %s WHERE %s\n", table.NameDB, whereClause)
		rows, err = DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s", table.NameDB, whereClause))
	}
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result, _, err := convertRowsToObject(value, rows)
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func addScenarioFields(values *[]string, fNames *[]string, scenarioVersionId string) {
	*values = append(*values, fmt.Sprintf("'%s'", scenarioVersionId))
	*fNames = append(*fNames, "ScenarioVersionId")
	*values = append(*values, "false")
	*fNames = append(*fNames, "IsDeleted")
}

func createWhereClauseAnd(values *[]string, fNames *[]string) string {
	var fields []string
	for i, fname := range *fNames {
		fields = append(fields, fmt.Sprintf("%s = %s", fname, (*values)[i]))
	}
	return strings.Join(fields, " AND ")
}

func createWhereClauseOr(values *[]string, fNames *[]string) string {
	var fields []string
	for i, fname := range *fNames {
		fields = append(fields, fmt.Sprintf("%s = %s", fname, (*values)[i]))
	}
	return strings.Join(fields, " OR ")
}

func convertRowsToObject[T any](value *T, rows *sql.Rows) ([]T, []string, error) {
	var result []T
	var resultScenario []string
	table := CreateTableFromStruct(value)

	for rows.Next() {
		var scIndex int
		s := reflect.ValueOf(value).Elem()
		numCols := s.NumField()
		if table.IsScenarioBased {
			numCols++
		}
		columns := make([]interface{}, numCols)
		columnIndex := 0
		for i := 0; i < s.NumField(); i++ {
			field := s.Field(i)
			//fmt.Printf("%s %s %s\n", field.Kind(), s.Type().Field(i).Name, field)
			if field.Type().Kind() != reflect.Struct {
				columns[columnIndex] = field.Addr().Interface()
			} else if s.Type().Field(i).Name == "Scenario" {
				var ScenarioVersionId string
				var IsDeleted bool
				columns[columnIndex] = &ScenarioVersionId
				columnIndex++
				columns[columnIndex] = &IsDeleted
				scIndex = i
			}
			columnIndex++
		}
		err := rows.Scan(columns...)
		if err != nil {
			return nil, nil, err
		}
		result = append(result, s.Interface().(T))
		if table.IsScenarioBased {
			resultScenario = append(resultScenario, *(columns[scIndex].(*string)))
		}
	}

	if table.IsScenarioBased {
		return result, resultScenario, nil
	} else {
		return result, nil, nil
	}

}
