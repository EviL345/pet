package httperrors

import (
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

const (
	ErrBadRequest         = "Bad request"
	ErrEmailAlreadyExists = "User with given email already exists"
	ErrNoSuchUser         = "User not found"
	ErrWrongCredentials   = "Wrong Credentials"
	ErrNotFound           = "Not Found"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrBadQueryParams     = "Invalid query params"
)

var (
	BadRequest            = errors.New("Bad request")
	WrongCredentials      = errors.New("Wrong Credentials")
	NotFound              = errors.New("Not Found")
	Unauthorized          = errors.New("Unauthorized")
	Forbidden             = errors.New("Forbidden")
	PermissionDenied      = errors.New("Permission Denied")
	ExpiredCSRFError      = errors.New("Expired CSRF token")
	WrongCSRFToken        = errors.New("Wrong CSRF token")
	CSRFNotPresented      = errors.New("CSRF not presented")
	NotRequiredFields     = errors.New("No such required fields")
	BadQueryParams        = errors.New("Invalid query params")
	InternalServerError   = errors.New("Internal Server Error")
	RequestTimeoutError   = errors.New("Request Timeout")
	ExistsEmailError      = errors.New("User with given email already exists")
	InvalidJWTToken       = errors.New("Invalid JWT token")
	InvalidJWTClaims      = errors.New("Invalid JWT claims")
	NotAllowedImageHeader = errors.New("Not allowed image header")
	NoCookie              = errors.New("not found cookie header")
)

type RestError interface {
	Status() int
	Error() string
	Causes() interface{}
}

type restErr struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() interface{} {
	return e.ErrCauses
}

func NewRestError(status int, err string, causes interface{}) RestError {
	return restErr{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func ErrorStatus(err error) int {
	if resterr, ok := err.(RestError); ok {
		return resterr.Status()
	}

	return http.StatusInternalServerError
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, ErrorStatus(err))

	render.JSON(w, r, err)
}
