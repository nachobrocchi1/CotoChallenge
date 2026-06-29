package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

// Encode writes any type serialized as JSON directly into the response.
func Encode[T any](w http.ResponseWriter, statusCode int, data T) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error encoding json response: %v", err)
	}
}

// ErrorResponse standard error response for the client
type ErrorResponse struct {
	Message string `json:"message"`
	Fields  any    `json:"fields,omitempty"`
}

// EncodeError sends a standardized error response back to the client.
func EncodeError(w http.ResponseWriter, statusCode int, err error) {
	Encode(w, statusCode, ErrorResponse{
		Message: err.Error(),
	})
}

// Decode parses and deserializes the request body stream into the target generic type.
// it uses a maxBytes limit to protect memory in case of larger requests.
func Decode[T any](body io.ReadCloser, maxBytes int64) (T, error) {
	defer body.Close()

	var target T

	// reads a safe amount of bytes
	limitedReader := io.LimitReader(body, maxBytes)

	dec := json.NewDecoder(limitedReader)
	dec.DisallowUnknownFields() // don't allow sending extra fields

	if err := dec.Decode(&target); err != nil {
		if err == io.EOF {
			return target, errors.New("request body cannot be empty")
		}
		return target, err
	}

	return target, nil
}

// EncodeNoContent set headers and the status code without writing a body stream.
func EncodeNoContent(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

type ValidationError struct {
	Fields map[string]string `json:"fields"`
}

func (e ValidationError) Error() string {
	var sb strings.Builder
	sb.WriteString("validation failed: ")
	for field, msg := range e.Fields {
		sb.WriteString(field)
		sb.WriteString(" (")
		sb.WriteString(msg)
		sb.WriteString("); ")
	}
	return sb.String()
}
