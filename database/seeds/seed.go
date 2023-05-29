package seeds

import (
	"log"
	"os"

	"github.com/jinzhu/gorm" // _ "github.com/jinzhu/gorm/dialects/postgres"

	"raptor/database"
)

// DB database instance
var DB *gorm.DB

// Init ...
func Init(config database.DBConfig) {
	var e error

	DB, e = database.New(config)

	if e != nil {
		log.Fatalln(e)

		os.Exit(1)
	}
}

// Seed ...
func Seed() {
	var err error

	// Create new sum type
	err = createSeed()
	if err != nil {
		log.Printf("Error seeding sum:%+v", err)
	}

}
