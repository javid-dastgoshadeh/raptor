package database

// DuplicateErr ...
type DuplicateErr struct {
	error
	Err map[string]interface{}
}

// RecordNotFoundErr ...
type RecordNotFoundErr struct {
	error
	Err interface{}
}

// UnableToCreateErr ...
type UnableToCreateErr struct {
	error
	Err interface{}
}
