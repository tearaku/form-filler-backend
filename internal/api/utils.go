package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"teacup1592/form-filler/internal/schoolForm"
)

/* Code from the following blog:
https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
*/

type RequestBody struct {
	UserId int    `json:"userId"`
	Secret string `json:"secret"`
}

// Maximum of 1MB of request body
const MAX_REQSIZE = 1048576

func validateRequest(w http.ResponseWriter, r *http.Request, service *schoolForm.Service) error {
	// Limit request size to prevent potentially malicious data
	r.Body = http.MaxBytesReader(w, r.Body, MAX_REQSIZE)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var rb RequestBody

	// Check if request format itself is ok
	err := dec.Decode(&rb)
	if err != nil {
		log.Println(err.Error())
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			msg := fmt.Sprintf("Request body contains malformed JSON @ position %d", syntaxErr.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		// Sometimes Decode() returns io.ErrUnexpectedEOF for JSON syntax err
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains malformed JSON"
			http.Error(w, msg, http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeErr):
			msg := fmt.Sprintf("Request body contains invalid value for %q field @ position %d", unmarshalTypeErr.Field, unmarshalTypeErr.Offset)
			http.Error(w, msg, http.StatusBadRequest)
		// Includes unknown fields, no sentinel error right now
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fName)
			http.Error(w, msg, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)
		// No sentinel error right now
		case err.Error() == "http: request body too large":
			msg := "Request body must be within 1MB"
			http.Error(w, msg, http.StatusRequestEntityTooLarge)
		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return err
	}

	// Check for origin of request
	if rb.Secret != os.Getenv("BACKEND_SECRET") {
		msg := "Invalid request"
		http.Error(w, msg, http.StatusInternalServerError)
		return errors.New("invalid request")
	}

	// Check for session
	if err := service.CheckSession(r.Context(), rb.UserId); err != nil {
		msg := "Invalid session"
		http.Error(w, msg, http.StatusBadRequest)
		return err
	}

	return nil
}
