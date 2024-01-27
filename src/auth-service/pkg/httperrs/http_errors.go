package httperrs

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	ErrUsernameOrPasswordInvalid = "email or password is incorrect"
	PasswordInvalid              = "invalid password"
	ErrEmailExisted              = "email already exists"
	CanNotSaveToStorage          = "an error occurred while storing data"
	PassAndConfirmPassNotMatch   = "password does not match the confirmed password"
	HasErrTryAgain               = "an error has occurred, please try again later!"
	ErrBody                      = "error body"
)

var (
	InternalServerError = errors.New("Internal Server Error")
	Unauthorized        = errors.New("Unauthorized")
	InvalidJWTToken     = errors.New("Invalid JWT token")
	InvalidJWTClaims    = errors.New("Invalid JWT claims")
)

type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error  Error() interface method
func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

// Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// New Unauthorized Error
func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		ErrCauses: causes,
	}
}
