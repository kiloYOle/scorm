package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

func AutoMigration(values ...interface{}) {
	fmt.Println("Starting auto-migration")
	fmt.Printf("%d objects to create\n", len(values))
	/*field1 := schema.Field{Name: "ID", NameDB: "ID", Type: "INT", PrimaryKey: true}
	fields := []*schema.Field{&field1}
	field2 := schema.Field{Name: "Other", NameDB: "Other", Type: "INT", PrimaryKey: false}
	fields = append(fields, &field2)
	table := schema.Table{Name: "Test", NameDB: "TestDB", Fields: fields}*/
	createScenario := false
	for _, value := range values {
		table := CreateTableFromStruct(value)
		if tableExists(DB, table) {
			fmt.Printf("Table %s exists\n", table.Name)
		} else {
			fmt.Printf("Creating table %s\n", table.Name)
			createTable(DB, table)
		}
		if table.IsScenarioBased {
			createScenario = true
		}
	}
	if createScenario {
		CreateScenarioTables()

	}
	fmt.Println("Auto-migration completed")
}

func CreateScenarioTables() {
	table := CreateTableFromStruct(&ScenarioTable{})
	createTable(DB, table)
	table = CreateTableFromStruct(&ScenarioVersionTable{})
	createTable(DB, table)
	fmt.Println("Scenario tables created")
}

func tableExists(db *sql.DB, table Table) bool {
	rows, err := db.Query("SELECT count(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ? AND table_type = ?", DBName, table.NameDB, "BASE TABLE")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tableCount := 0
	if rows.Next() {
		if err := rows.Scan(&tableCount); err != nil {
			panic(err)
		}
	}

	return tableCount != 0
}

func createTable(db *sql.DB, table Table) {
	createString := "CREATE TABLE " + table.NameDB + " "
	fieldString := createFieldText(table.Fields)
	statement := createString + " ( " + fieldString + " )"
	fmt.Println(statement)
	_, err := db.Exec(statement)
	if err != nil {
		panic(err)
	}

}

func createFieldText(fields []*Field) string {
	var fieldStrings []string
	var primaryKeyFields []string
	for _, f := range fields {
		fieldString := f.Name + " " + f.Type
		fieldStrings = append(fieldStrings, fieldString)
		if f.PrimaryKey {
			primaryKeyFields = append(primaryKeyFields, f.NameDB)
		}
	}
	result := strings.Join(fieldStrings, ", ") + ", PRIMARY KEY(" + strings.Join(primaryKeyFields, ",") + ")"

	return result
}
