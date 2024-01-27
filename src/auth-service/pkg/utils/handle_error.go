package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

func ErrorsBindParamOrBody(err error) []ErrorResponse {
	var ve validator.ValidationErrors
	var out []ErrorResponse
	if errors.As(err, &ve) {
		for _, fe := range ve {
			out = append(out, ErrorResponse{Detail: fe.Error(), Message: ValidationErrorToText(fe)})
		}
	}
	return out
}
