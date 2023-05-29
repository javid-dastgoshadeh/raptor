package migrations

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"

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

// Migrate ...
func Migrate() {
	sample(DB)
	fmt.Println("Migration Done!")
}

// RunScripts , this script contains most important functions
// and extensions that app requires
func RunScripts() {

	// postgres
	dat, err := ioutil.ReadFile("database/scripts/postgres")

	if err != nil {
		panic(err)
	}

	psQuery := string(dat)
	DB.Exec(psQuery)
	log.Println("Script run successfully")
}
