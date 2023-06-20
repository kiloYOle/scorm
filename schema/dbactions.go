package schema

import (
	"fmt"
	"reflect"
	"strings"
)

func Insert(value interface{}, scenario string) {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	if table.IsScenarioBased {
		addScenarioFields(&values, &fNames, scenario)
	}
	insertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(fNames, ","), strings.Join(values, ","))
	fmt.Println(insertString)
	_, err := DB.Exec(insertString)
	if err != nil {
		fmt.Println("Insert error", err)
	}
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
			columns[i] = field.Addr().Interface()
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
	fmt.Println("len before", len(*values))
	*values = append(*values, fmt.Sprintf("'%s'", scenarioId))
	*fNames = append(*fNames, "scenarioId")
	fmt.Println("len after", len(*values))
}
