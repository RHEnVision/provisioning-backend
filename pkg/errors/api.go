// Copyright Red Hat

package errors

import (
	"encoding/json"
	"net/http"
)

// APIError defines a type for all errors returned by the IDP Config service
type APIError struct {
	Code   string `json:"Code"`
	Status int    `json:"Status"`
	Title  string `json:"Title"`
}

// Error gets a error string from an APIError
func (e *APIError) Error() string { return e.Title }

// InternalServerError defines a generic error for the IDP Config service
type InternalServerError struct {
	APIError
}

// NewInternalServerError creates a new InternalServerError
func NewInternalServerError(message string) *InternalServerError {
	err := new(InternalServerError)
	err.Code = "ERROR"
	err.Title = message
	err.Status = http.StatusInternalServerError
	return err
}

// Conflict defines a 409 error for the IDP Config service (Violation of Unique Constraint for name within an account)
type Conflict struct {
	APIError
}

// NewConflict creates a new Conflict
func NewConflict(message string) *Conflict {
	err := new(Conflict)
	err.Code = "ERROR"
	err.Title = message
	err.Status = http.StatusConflict
	return err
}

// Forbidden defines a 403 error
type Forbidden struct {
	APIError
}

// NewForbidden creates a new Conflict
func NewForbidden(message string) *Forbidden {
	err := new(Forbidden)
	err.Code = "ERROR"
	err.Title = message
	err.Status = http.StatusForbidden
	return err
}

// BadRequest defines a error when the client's input generates an error
type BadRequest struct {
	APIError
}

// NewBadRequest creates a new BadRequest
func NewBadRequest(message string) *BadRequest {
	err := new(BadRequest)
	err.Code = "BAD_REQUEST"
	err.Title = message
	err.Status = http.StatusBadRequest
	return err
}

// NotFound defines a error for whenever an entity is not found in the database
type NotFound struct {
	APIError
}

// NewNotFound creates a new NotFound
func NewNotFound(message string) *NotFound {
	err := new(NotFound)
	err.Code = "NOT_FOUND"
	err.Title = message
	err.Status = http.StatusNotFound
	return err
}

func RespondWithBadRequest (message string, w http.ResponseWriter) {
	err := NewBadRequest(message)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(&err)
}

func RespondWithInternalServerError (message string, w http.ResponseWriter) {
	err := NewInternalServerError(message)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(&err)
}

func RespondWithConflict (message string, w http.ResponseWriter) {
	err := NewConflict(message)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(&err)
}

func RespondWithForbidden (message string, w http.ResponseWriter) {
	err := NewForbidden(message)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(&err)
}

func RespondWithNotFound (message string, w http.ResponseWriter) {
	err := NewNotFound(message)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(&err)
}
