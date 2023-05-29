package seeds

import (
	"raptor/database"
)

// CreateSum ...
func createSeed() error {

	err := DB.Create("").Error
	if err != nil {
		return database.UnableToCreateErr{Err: err.Error()}
	}

	err = DB.Save("").Error
	if err != nil {
		return database.UnableToCreateErr{Err: err.Error()}
	}
	return nil
}
