// Implementation taken from gilcrest/diy-go-api
// with very slight modifications

// I guess I just need to put this here...?
// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Package errs defines the error handling used by all Upspin software.
package errs

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

type Error struct {
	// User ID of the user attempting the operation
	// 0 if unobtainable
	User User
	// Class of error (ex: permission failure),
	// or, "Other" if class is unknown or irrelevant
	Kind Kind
	// The parameters related to the error
	Param Parameter
	// Human-readable, short description of the error
	Code Code
	// Underlying error triggering this error, if any
	Err error
}

func (e *Error) isEmpty() bool {
	return e.User == 0 && e.Kind == 0 && e.Param == "" && e.Code == "" && e.Err == nil
}

// Allows for unwrapping errors with errors.As
func (e Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// Defines user id of user attempting the operation
// 0 if unobtainable
type User uint32

// Defines kind of error, mostly for use by systems like FUSE
// that must act differently depending on error
type Kind uint8

// Parameter related to the error
type Parameter string

// Human-readable, short description of the error
type Code string

// Kinds of errors
const (
	Other           Kind = iota // Unclassified error; this value is NOT printed in err msg
	Invalid                     // Invalid operation for this type of item
	IO                          // External I/O error, like network failure
	Exist                       // Item already exists
	NotExist                    // Item does NOT exist
	Private                     // Information withheld
	Internal                    // Internal error, or inconsistency
	Database                    // Error from DB
	Unauthenticated             // Returns 401, empty body
	Unauthorized                // Returns 403, empty body
	InvalidRequest
	Validation
	Unanticipated
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other_error"
	case Invalid:
		return "invalid_operation"
	case IO:
		return "I/O_error"
	case Exist:
		return "item_already_exists"
	case NotExist:
		return "item_does_not_exist"
	case Private:
		return "information_withheld"
	case Internal:
		return "internal_error"
	case Database:
		return "database_error"
	case Validation:
		return "input_validation_error"
	case Unanticipated:
		return "unanticipated_error"
	case InvalidRequest:
		return "invalid_request_error"
	case Unauthenticated:
		return "unauthenticated_request"
	case Unauthorized:
		return "unauthorized_request"
	}
	return "unknown_error_kind"
}

// Builds error from arguments
// Must have at least 1 argument, else it panics
// Type of the arguments determines its meaning (AKA checked by type)
// If more than 1 argument of a given type is given, only the last takes effect
//
// The types are:
//
//	UserName
//		The username of the user attempting the operation.
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.New.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case User:
			e.User = arg
		case string:
			e.Err = errors.New(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			e.Err = arg
		case error:
			// TODO: removed stack trace, do I need it?
			e.Err = arg
		case Code:
			e.Code = arg
		case Parameter:
			e.Param = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			return fmt.Errorf("errors.E: bad call from %s:%d: %v, unknown type %T, value %v in error call", file, line, args, arg, arg)
		}
	}

	// Get the underlying error that triggered this
	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}
	// If error has Kind unset or Other, pull up inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}

	if prev.Code == e.Code {
		prev.Code = ""
	}
	// If this error has Code == "", pull up the inner one.
	if e.Code == "" {
		e.Code = prev.Code
		prev.Code = ""
	}

	if prev.Param == e.Param {
		prev.Param = ""
	}
	// If this error has Param == "", pull up the inner one.
	if e.Param == "" {
		e.Param = prev.Param
		prev.Param = ""
	}

	return e
}

// Match compares its two error arguments. It can be used to check
// for expected errors in tests. Both arguments must have underlying
// type *Error or Match will return false. Otherwise it returns true
// if every non-zero element of the first error is equal to the
// corresponding element of the second.
// If the Err field is a *Error, Match recurs on that field;
// otherwise it compares the strings returned by the Error methods.
// Elements that are in the second argument but not present in
// the first are ignored.
//
// For example,
//
// Match(errors.E(upspin.UserName("joe@schmoe.com"), errors.Permission), err)
// tests whether err is an Error with Kind=Permission and User=joe@schmoe.com.
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.User != 0 && e2.User != e1.User {
		return false
	}
	if e1.Kind != Other && e2.Kind != e1.Kind {
		return false
	}
	if e1.Param != "" && e2.Param != e1.Param {
		return false
	}
	if e1.Code != "" && e2.Code != e1.Code {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

// KindIs reports whether err is an *Error of the given Kind.
// If err is nil then KindIs returns false.
func KindIs(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return KindIs(kind, e.Err)
	}
	return false
}
