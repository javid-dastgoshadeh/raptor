package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	// Include postgres dilect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var _db *gorm.DB

// DBConfig ...
type DBConfig struct {
	Engine   string `json:"engine"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	SslMode  string `json:"ssl_mode"`
	Log      bool   `json:"log"`
}

// New ...
func New(cnf DBConfig) (*gorm.DB, error) {
	// Check if there is an older instance and return it,
	// else creates a new instance and stores it in _db object
	if _db != nil {
		return _db, nil
	}

	var err error

	switch cnf.Engine {
	case "postgres":
		_db, err = gorm.Open(cnf.Engine, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cnf.Host,
			cnf.Port,
			cnf.Username,
			cnf.DbName,
			cnf.Password,
			cnf.SslMode))
	case "mysql":
	case "mssql":
	default:
		_db, err = gorm.Open(cnf.Engine, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cnf.Host,
			cnf.Port,
			cnf.Username,
			cnf.DbName,
			cnf.Password,
			cnf.SslMode))
	}

	if err != nil {
		return nil, err
	}

	_db.LogMode(cnf.Log)

	return _db, err
}

// GetInstance ...
func GetInstance() *gorm.DB {
	return _db
}
