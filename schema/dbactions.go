package schema

import (
	"fmt"
	"reflect"
	"strings"
)

func Insert(value interface{}, scenarioId string) error {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	if table.IsScenarioBased {
		scenario, err := Find(&ScenarioTable{ScenarioId: scenarioId}, "")
		if len(scenario) > 0 {
			addScenarioFields(&values, &fNames, scenarioId)
		} else if err != nil {
			return err
		} else {
			fmt.Println("Scenario not found, record not inserted")
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
	pkValues, pkNames := CreatePKFieldValuesAndNames(table, value)
	whereClause := createWhereClauseAnd(&pkValues, &pkNames)
	fmt.Printf("SELECT * FROM %s WHERE %s\n", table.NameDB, whereClause)
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s WHERE %s", table.NameDB, whereClause))
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer rows.Close()

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
		fmt.Println(s)
		result = append(result, s.Interface().(T))
	}
	fmt.Println(result)
	return result, err
}

func FindAll[T any](value *T, scenario string) ([]T, error) {
	var result []T
	table := CreateTableFromStruct(value)
	//values, fNames := CreateFieldValuesAndNames(table, value)
	//pkValues, pkNames := CreatePKFieldValuesAndNames(table, value)
	fmt.Printf("SELECT * FROM %s\n", table.NameDB)
	rows, err := DB.Query(fmt.Sprintf("SELECT * FROM %s", table.NameDB))
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer rows.Close()

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
		fmt.Println(s)
		result = append(result, s.Interface().(T))
	}
	fmt.Println(result)
	return result, err
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
