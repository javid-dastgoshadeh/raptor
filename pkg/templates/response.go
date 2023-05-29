package templates

import (
	"net/http"
	"raptor/pkg/failure"

	"raptor/config"
)

// ResponseTemplate standard template for http responses
type ResponseTemplate struct {
	Status  string          `json:"status,omitempty"`
	Code    failure.IntCode `json:"code,omitempty"`
	Message interface{}     `json:"message,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    interface{}     `json:"meta,omitempty"`
	Links   interface{}     `json:"links,omitempty"`
}

// PaginateTemplate ...
type PaginateTemplate struct {
	Total int     `json:"total"`
	Pages int     `json:"pages"`
	Limit int     `json:"limit"`
	Page  int     `json:"page"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// GetWithCode return template with considering error code
func GetWithCode(code int, err error) (template *ResponseTemplate) {
	msg := GetMessage(code, err)

	switch code {
	case http.StatusInternalServerError:
		template = InternalServerError(msg)
	case http.StatusBadRequest:
		template = BadRequest(msg)
	case http.StatusForbidden:
		template = Forbidden(msg)
	case http.StatusNotFound:
		template = NotFound(msg)
	case http.StatusUnprocessableEntity:
		template = UnprocessableEntity(msg)
	case http.StatusUnauthorized:
		template = Unauthorized(msg)
	case http.StatusMethodNotAllowed:
		template = MethodNotAllowed(msg)
	default:
		template = InternalServerError(msg)
	}

	return
}

// GetMessage return error message
func GetMessage(code int, err error) (msg string) {
	if env.GetBool("debug") {
		msg = err.Error()

		return
	}

	switch code {
	case http.StatusInternalServerError:
		msg = "Internal server error"
	case http.StatusBadRequest:
		msg = "Bad request"
	case http.StatusForbidden:
		msg = "Forbidden"
	case http.StatusNotFound:
		msg = "Not found"
	case http.StatusUnprocessableEntity:
		msg = "Unprocessable entity"
	case http.StatusUnauthorized:
		msg = "Unauthorized"
	case http.StatusMethodNotAllowed:
		msg = "MethodNotAllowed"
	case http.StatusOK:
		msg = "Ok"
	default:
		msg = "Internal server error"
	}

	return
}

// BadRequest ...
func BadRequest(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusBadRequest,
		Status: "fail",
		//Message: "Bad request",
		Data: map[string]interface{}{
			"message": msg,
		},
	}
}

// InternalServerError ...
func InternalServerError(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:   http.StatusInternalServerError,
		Status: "error",
		//Message: "Internal server error",
		Message: map[string]interface{}{
			"message": msg,
		},
	}
}

// MobileAppInternalServerError ...
func MobileAppInternalServerError(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:   http.StatusInternalServerError,
		Status: "error",
		Message: map[string]interface{}{
			"message": msg,
		},
	}
}

// NotFound ...
func NotFound(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusNotFound,
		Status: "fail",
		//Message: "Not found",
		Data: map[string]interface{}{
			"message": msg,
		},
	}
}

// UnprocessableEntity ...
func UnprocessableEntity(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		Code:    http.StatusUnprocessableEntity,
		Status:  "UNPROCESSABLE_ENTITY",
		Message: "Unprocessable entity",
		Data: map[string]interface{}{
			"message": msg,
		},
	}
}

// Unauthorized ...
func Unauthorized(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusUnauthorized,
		Status: "fail",
		//Message: "Unauthorized",
		Data: map[string]interface{}{
			//"message": msg,
			"message": "unauthorized request",
		},
	}
}

// Forbidden ...
func Forbidden(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusForbidden,
		Status: "fail",
		//Message: "Forbidden",
		Data: map[string]interface{}{
			"message": msg,
		},
	}
}

// MethodNotAllowed ...
func MethodNotAllowed(msg interface{}) *ResponseTemplate {

	return &ResponseTemplate{
		//Code:    http.StatusMethodNotAllowed,
		Status: "fail",
		//Message: "Method not allowed",
		Data: map[string]interface{}{
			//"message": msg,
			"message": "method not allowed",
		},
	}
}

// Ok ...
func Ok(data interface{}, meta interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusOK,
		Status: "success",
		//Message: "Ok",
		Data: data,
		//Meta: meta,
	}
}

// MobileAppRegisterResponse ...
func MobileAppRegisterResponse(msg interface{}) *ResponseTemplate {
	return &ResponseTemplate{
		//Code:    http.StatusOK,
		Status: "success",
		//Message: "Ok",
		Data: msg,
	}
}
