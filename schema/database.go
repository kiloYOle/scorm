package schema

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB
var DBName string

func OpenDBConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBName = os.Getenv("DBNAME")
	// Capture connection properties.
	cfg := mysql.Config{
		User:              os.Getenv("DBUSER"),
		Passwd:            os.Getenv("DBPASS"),
		Net:               "tcp",
		Addr:              os.Getenv("DBADDR"),
		DBName:            os.Getenv("DBNAME"),
		InterpolateParams: true,
	}
	// Get a database handle.
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := DB.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Printf("Connected to database at %s\n", cfg.Addr)
}
