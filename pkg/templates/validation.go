package templates

import (
	"fmt"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

// FieldError , format for each errors off a field
type FieldError struct {
	Code string
	Err  string
}

// ValidationErrTemplates ...
var ValidationErrTemplates = map[string]string{
	"required": "Filed ${field} is required",
	"min":      "Minimum length(value) for ${field} is ${param}",
	"max":      "Maximum length(value) for ${field} is ${param}",
	"oneof":    "Filed ${field} must be one of (${param})",
	"email":    "Filed ${field} must be a email",
	"unique":   "Filed ${field} must be unique in list",
	"number":   "Filed ${field} must be a number",
}

// ValidationErrCodes ...
var ValidationErrCodes = map[string]string{
	"required": "Required",
	"min":      "Min",
	"max":      "Max",
	"oneof":    "Oneof",
	"email":    "Email",
	"unique":   "Unique",
	"number":   "Number",
}

// ParseErrors parse errors to special format
func ParseErrors(validationErrors *validator.ValidationErrors) map[string][]FieldError {
	var (
		field,
		tag,
		errStr string
		errMap = make(map[string][]FieldError, len(*validationErrors))
	)

	for _, fieldError := range *validationErrors {
		field = fieldError.Field()
		tag = fieldError.Tag()

		errStr = GetError(fieldError)

		fieldErr := FieldError{Code: GetErrCode(tag), Err: errStr}

		errMap[field] = append(errMap[field], fieldErr)
	}

	return errMap
}

// ParseDuplicateErrors ...
func ParseDuplicateErrors(duplicateError map[string]interface{}) map[string]FieldError {

	err := fmt.Sprintf("a record with given %s is already exists", duplicateError["field"])
	errMap := map[string]FieldError{
		duplicateError["field"].(string): FieldError{
			Code: duplicateError["field"].(string),
			Err:  err,
		},
	}

	return errMap
}

// FillTemplate ...
func FillTemplate(template string, field, tag, param string) string {
	var re = regexp.MustCompile(`(\$\{field\})`)
	s := re.ReplaceAllString(template, field)

	re = regexp.MustCompile(`(\$\{tag\})`)
	s = re.ReplaceAllString(s, tag)

	re = regexp.MustCompile(`(\$\{param\})`)
	s = re.ReplaceAllString(s, param)

	return s
}

// GetError return error message relative to tag
func GetError(fieldError validator.FieldError) string {
	var (
		field = fieldError.Field()
		// value = fieldError.Value().(string)
		tag   = fieldError.Tag()
		param = fieldError.Param()
	)

	if template, ok := ValidationErrTemplates[fieldError.Tag()]; ok {
		return FillTemplate(template, field, tag, param)
	}

	return "Unknown error"
}

// GetErrCode , return error code for tag
func GetErrCode(tag string) string {
	if val, ok := ValidationErrCodes[tag]; ok {
		return val
	}

	return tag
}
