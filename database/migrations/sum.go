package migrations

import (
	"github.com/jinzhu/gorm"
)

// Sum ...
func sample(db *gorm.DB) {
	db.AutoMigrate()
}
