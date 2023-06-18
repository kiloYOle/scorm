package main

import (
	"fmt"

	"github.com/kiloYOle/scorm/migrator"
	"github.com/kiloYOle/scorm/schema"
)

func main() {
	fmt.Println("Hello, world.")
	schema.OpenDBConnection()
	migrator.AutoMigration()
}
