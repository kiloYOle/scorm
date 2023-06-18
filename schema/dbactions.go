package schema

import (
	"fmt"
	"strings"
)

func Insert(value interface{}) {
	table := CreateTableFromStruct(value)
	values, fNames := CreateFieldValuesAndNames(table, value)
	insertString := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table.Name, strings.Join(fNames, ","), strings.Join(values, ","))
	fmt.Println(insertString)
	_, err := DB.Exec(insertString)
	if err != nil {
		fmt.Println(err)
	}
}

func Update(value interface{}) {
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

func Delete(value interface{}) {
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
