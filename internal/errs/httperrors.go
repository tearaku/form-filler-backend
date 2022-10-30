// I guess I can just put this here...?
// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

// ErrResponse is used as the Response Body
type ErrorResponse struct {
	Error ServiceError `json:"error"`
}

// ServiceError has fields for Service errors. All fields with no data will
// be omitted
type ServiceError struct {
	Kind    string `json:"kind,omitempty"`
	Code    string `json:"code,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

// HTTPErrorResponse takes a writer, error and a logger, performs a
// type switch to determine if the type is an Error (which meets
// the Error interface as defined in this package), then sends the
// Error as a response to the client. If the type does not meet the
// Error interface as defined in this package, then a proper error
// is still formed and sent to the client, however, the Kind and
// Code will be Unanticipated. Logging of error is also done using
// https://github.com/rs/zerolog
// TODO: replace built-in logger w/ zerolog
func HTTPErrorResponse(w http.ResponseWriter, lgr zerolog.Logger, err error) {
	if err == nil {
		nilErrResponse(w, lgr)
		return
	}

	var e *Error
	if errors.As(err, &e) {
		switch e.Kind {
		case Unauthorized:
			unauthorizedErrResponse(w, lgr, e)
			return
		default:
			typicalErrResponse(w, lgr, e)
			return
		}
	}
	unknownErrResponse(w, lgr, err)
}

// Replies to the request with the specified error
// message and HTTP code. It does not otherwise end the request; the
// caller should ensure no further writes are done to w
func typicalErrResponse(w http.ResponseWriter, lgr zerolog.Logger, e *Error) {
	statusCode := httpErrStatusCode(e.Kind)

	// Simply send status code to client if empty error occurs (though it should NOT)
	if e.isEmpty() {
		lgr.Error().Stack().Int("http_statuscode", statusCode).Msg("Empty error")
		// TODO: check this again...
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Typical error flow
	lgr.Error().Stack().Err(e.Err).
		Int("http_statuscode", statusCode).
		Str("Kind", e.Kind.String()).
		Str("Parameter", string(e.Param)).
		Str("Code", string(e.Code)).
		Msg("Error response: sent.")

	eR := newErrResponse(e)
	errJSON, _ := json.Marshal(eR)
	eJ := string(errJSON)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	// Write response body (json)
	fmt.Fprintln(w, eJ)
}

func newErrResponse(err *Error) ErrorResponse {
	msg := "internal server error - please contact dev"

	switch err.Kind {
	case Internal, Database:
		return ErrorResponse{
			Error: ServiceError{
				Kind:    Internal.String(),
				Message: msg,
			},
		}
	default:
		return ErrorResponse{
			Error: ServiceError{
				Kind:    err.Kind.String(),
				Code:    string(err.Code),
				Param:   string(err.Param),
				Message: err.Error(),
			},
		}
	}
}

// Responds w/ 403 (Forbidden) & empty response body
func unauthorizedErrResponse(w http.ResponseWriter, lgr zerolog.Logger, err *Error) {
	lgr.Error().Stack().Err(err.Err).
		Int("http_statuscode", http.StatusForbidden).
		Msg("Unauthorized request")
	w.WriteHeader(http.StatusInternalServerError)
}

func unknownErrResponse(w http.ResponseWriter, lgr zerolog.Logger, err error) {
	eR := ErrorResponse{
		Error: ServiceError{
			Kind:    Unanticipated.String(),
			Code:    "Unanticipated",
			Message: "Unexpected error: please contact the dev",
		},
	}
	lgr.Error().Err(err).Msg("Unknown error")

	errJSON, _ := json.Marshal(eR)
	eJ := string(errJSON)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, eJ)
}

// Should not be triggered, just in case
func nilErrResponse(w http.ResponseWriter, lgr zerolog.Logger) {
	lgr.Error().Stack().
		Int("http_statuscode", http.StatusInternalServerError).
		Msg("nil Error: no response body sent")
	w.WriteHeader(http.StatusInternalServerError)
}

// Maps error Kind --> a HTTP status code
func httpErrStatusCode(k Kind) int {
	switch k {
	case Invalid, Exist, NotExist, Private, Validation, InvalidRequest:
		return http.StatusBadRequest
	// Zero value of Kind is Other, so if no Kind is present in the error,
	// Other is used. Errors should always have a Kind set, otherwise, a
	// 500 will be returned and no error message will be sent to the caller
	case Other, IO, Internal, Database, Unanticipated:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
