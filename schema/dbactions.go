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
		scenario, err := getScenario(scenarioId)
		if err == nil {
			addScenarioFields(&values, &fNames, scenario.ScenarioId)
		} else {
			fmt.Println(err)
			return nil
		}
	}
	insertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(fNames, ","), strings.Join(values, ","))
	fmt.Println("insert string")
	fmt.Println(insertString)
	_, err := DB.Exec(insertString)
	if err != nil {
		fmt.Println("Insert error", err)
	}
	return nil
}

func Update(value interface{}, scenario string) {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	pkValues, pkNames := CreatePKFieldValuesAndNames(table, value)
	var setFieldString []string
	for i, fName := range fNames {
		setFieldString = append(setFieldString, fName+"="+values[i])
	}
	var whereString []string
	for i, fName := range pkNames {
		whereString = append(whereString, fName+"="+pkValues[i])
	}
	updateString := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table.Name, strings.Join(setFieldString, ","), strings.Join(whereString, " AND "))
	fmt.Println(updateString)
	_, err := DB.Exec(updateString)
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

func Find[T any](value *T, scenario string) ([]T, error) {
	var result []T
	table := CreateTableFromStruct(value)
	//values, fNames := CreateFieldValuesAndNames(table, value)
	values, names := CreatePKFieldValuesAndNames(table, value)
	whereClause := createWhereClauseAnd(&values, &names)
	if table.IsScenarioBased {
		scenarios, _ := getScenarioAndParents(scenario)
		scValues, scNames := createScenarioFieldValuesAndNames(scenarios)
		whereClause_sc := createWhereClauseOr(&scValues, &scNames)
		whereClause = fmt.Sprintf("%s AND (%s)", whereClause, whereClause_sc)
	}
	fmt.Printf("SELECT * FROM %s WHERE %s\n", table.NameDB, whereClause)
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s", table.NameDB, whereClause))
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer rows.Close()

	return convertRowsToObject(value, rows)
}

func FindAll[T any](value *T, scenario string) ([]T, error) {
	table := CreateTableFromStruct(value)
	var whereClause string
	if table.IsScenarioBased {
		scenarios, _ := getScenarioAndParents(scenario)
		scValues, scNames := createScenarioFieldValuesAndNames(scenarios)
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

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	return convertRowsToObject(value, rows)
}

func addScenarioFields(values *[]string, fNames *[]string, scenarioId string) {
	*values = append(*values, fmt.Sprintf("'%s'", scenarioId))
	*fNames = append(*fNames, "scenarioId")
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

func convertRowsToObject[T any](value *T, rows *sql.Rows) ([]T, error) {
	var result []T

	for rows.Next() {
		s := reflect.ValueOf(value).Elem()
		numCols := s.NumField()
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			field := s.Field(i)
			//fmt.Printf("%s %s %s\n", field.Kind(), s.Type().Field(i).Name, field)
			if field.Type().Kind() != reflect.Struct {
				columns[i] = field.Addr().Interface()
			} else if s.Type().Field(i).Name == "Scenario" {
				var scenarioId string
				columns[i] = &scenarioId
			}
		}
		err := rows.Scan(columns...)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, s.Interface().(T))
	}

	return result, nil
}
