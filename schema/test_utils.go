package schema

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type InsertTest1 struct {
	Name string `scorm:"pk"`
	Flag bool
	Nr   int
}

type InsertTestScenario1 struct {
	Scenario
	Name string `scorm:"pk"`
	Flag bool
	Nr   int
}

func RunTestMain(m *testing.M) (code int, err error) {
	// need to define the relative path to godotenv while running tests
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBName = os.Getenv("TEST_DBNAME")
	// Capture connection properties.
	cfg := mysql.Config{
		User:              os.Getenv("TEST_DBUSER"),
		Passwd:            os.Getenv("TEST_DBPASS"),
		Net:               "tcp",
		Addr:              os.Getenv("TEST_DBADDR"),
		InterpolateParams: true,
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	DB.Exec(fmt.Sprintf("DROP DATABASE %s", DBName))
	fmt.Printf("CREATE DATABASE %s\n", DBName)
	_, err = DB.Exec(fmt.Sprintf("CREATE DATABASE %s;", DBName))
	if err != nil {
		panic(err)
	}
	_, err = DB.Exec("USE " + DBName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Test database created")
	// delete test database after all done
	defer func() {
		_, err := DB.Exec(fmt.Sprintf("DROP DATABASE %s", DBName))
		if err != nil {
			panic(err)
		}
		DB.Close()
		fmt.Println("Test database deleted")
	}()

	return m.Run(), nil
}

func OpenTestDB() {
	// need to define the relative path to godotenv while running tests
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBName = os.Getenv("TEST_DBNAME")
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("TEST_DBUSER"),
		Passwd: os.Getenv("TEST_DBPASS"),
		Net:    "tcp",
		Addr:   os.Getenv("TEST_DBADDR"),
		DBName: os.Getenv("TEST_DBNAME"),
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTestFieldListPointer() []*Field {
	field1 := Field{Name: "ID", NameDB: "ID", Type: "INT", PrimaryKey: true}
	fields := []*Field{&field1}
	return fields
}
