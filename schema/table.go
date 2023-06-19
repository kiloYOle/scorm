package schema

import (
	"fmt"
	"reflect"
)

type Table struct {
	Name            string
	NameDB          string
	IsScenarioBased bool
	Fields          []*Field
	FieldsByName    map[string]*Field
	FieldsByDBName  map[string]*Field
}

func CreateTableFromStruct(q interface{}) Table {
	fields := []*Field{}
	isScenarioBased := false

	v := reflect.ValueOf(q).Elem()
	//fmt.Printf("Table name %s\n", v.Type().Name())
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		//n := v.Type().Field(j).Name
		//t := f.Type().Name()
		tags, _ := v.Type().Field(j).Tag.Lookup("scorm")
		isPK := tags == "pk"
		//fmt.Printf("Field Name: %s Type: %s PK: %t\n", n, t, isPK)
		if f.Type().Name() == "Scenario" {
			field := Field{Name: "ScenarioId", NameDB: "ScenarioId", Type: GoTypesToDbTypes["string"], PrimaryKey: true, ScenarioField: true}
			fields = append(fields, &field)
			isScenarioBased = true
		} else {
			field := Field{Name: v.Type().Field(j).Name, NameDB: v.Type().Field(j).Name, Type: GoTypesToDbTypes[f.Type().Name()], PrimaryKey: isPK, ScenarioField: false}
			fields = append(fields, &field)
		}
	}

	fieldMapName := make(map[string]*Field)
	fieldMapDbName := make(map[string]*Field)
	for _, field := range fields {
		fieldMapName[field.Name] = field
		fieldMapDbName[field.NameDB] = field
	}

	table := Table{Name: v.Type().Name(), NameDB: v.Type().Name(), Fields: fields,
		FieldsByName: fieldMapName, FieldsByDBName: fieldMapDbName, IsScenarioBased: isScenarioBased}

	return table
}

func CreateFieldValuesAndNames(table Table, row interface{}) ([]string, []string) {
	values := []string{}
	fieldNames := []string{}

	v := reflect.ValueOf(row).Elem()
	for j := 0; j < v.NumField(); j++ {
		f := v.Type().Field(j)
		if f.Type.Kind() != reflect.Struct {
			values = append(values, ReflectValueToDBValue(v.Field(j)))
			fieldNames = append(fieldNames, table.FieldsByName[f.Name].NameDB)
		}
	}
	return values, fieldNames
}

func CreatePKFieldValuesAndNames(table Table, row interface{}) ([]string, []string) {
	values := []string{}
	fieldNames := []string{}

	v := reflect.ValueOf(row).Elem()
	for j := 0; j < v.NumField(); j++ {
		f := v.Type().Field(j)
		tags, _ := v.Type().Field(j).Tag.Lookup("scorm")
		isPK := tags == "pk"
		if isPK {
			values = append(values, ReflectValueToDBValue(v.Field(j)))
			fieldNames = append(fieldNames, table.FieldsByName[f.Name].NameDB)
		}
	}
	return values, fieldNames
}

func ReflectValueToDBValue(rValue reflect.Value) string {
	switch rValue.Kind() {
	case reflect.String:
		return fmt.Sprintf("'%s'", rValue.String())
	case reflect.Int:
		return fmt.Sprintf("%d", rValue.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", rValue.Bool())
	default:
		return ""
	}
}
