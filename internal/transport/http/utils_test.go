package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		data            any
		expectedStatus  int
		expectedBody    string
		expectedContent string
	}{
		{
			name:            "encodes struct as json",
			statusCode:      http.StatusOK,
			data:            map[string]string{"key": "value"},
			expectedStatus:  http.StatusOK,
			expectedBody:    `{"key":"value"}`,
			expectedContent: "application/json; charset=utf-8",
		},
		{
			name:            "encodes with status 201",
			statusCode:      http.StatusCreated,
			data:            map[string]int{"id": 1},
			expectedStatus:  http.StatusCreated,
			expectedBody:    `{"id":1}`,
			expectedContent: "application/json; charset=utf-8",
		},
		{
			name:            "encodes nil data",
			statusCode:      http.StatusOK,
			data:            nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    `null`,
			expectedContent: "application/json; charset=utf-8",
		},
		{
			name:            "encodes empty struct",
			statusCode:      http.StatusOK,
			data:            struct{}{},
			expectedStatus:  http.StatusOK,
			expectedBody:    `{}`,
			expectedContent: "application/json; charset=utf-8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Encode(w, tt.statusCode, tt.data)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != tt.expectedContent {
				t.Errorf("expected Content-Type %q, got %q", tt.expectedContent, contentType)
			}

			var got, expected any
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatalf("response body is not valid json: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.expectedBody), &expected); err != nil {
				t.Fatalf("expected body is not valid json: %v", err)
			}
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", expected) {
				t.Errorf("expected body %v, got %v", expected, got)
			}
		})
	}
}

type DummyStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecode(t *testing.T) {
	type args struct {
		body     io.ReadCloser
		maxBytes int64
	}

	tests := []struct {
		name    string
		args    args
		want    DummyStruct
		wantErr bool
		errSign string // Subcadena que esperamos encontrar en el mensaje de error
	}{
		{
			name: "Success - Valid JSON within byte limits",
			args: args{
				body:     io.NopCloser(strings.NewReader(`{"name":"John","age":30}`)),
				maxBytes: 1024,
			},
			want:    DummyStruct{Name: "John", Age: 30},
			wantErr: false,
		},
		{
			name: "Failure - Completely empty body (io.EOF)",
			args: args{
				body:     io.NopCloser(strings.NewReader("")),
				maxBytes: 1024,
			},
			want:    DummyStruct{},
			wantErr: true,
			errSign: "request body cannot be empty",
		},
		{
			name: "Failure - Malformed JSON syntax",
			args: args{
				body:     io.NopCloser(strings.NewReader(`{"name":"John",`)),
				maxBytes: 1024,
			},
			want:    DummyStruct{},
			wantErr: true,
			errSign: "unexpected EOF",
		},
		{
			name: "Failure - Unknown fields disallowed",
			args: args{
				body:     io.NopCloser(strings.NewReader(`{"name":"John","age":30,"extra":"field"}`)),
				maxBytes: 1024,
			},
			want:    DummyStruct{},
			wantErr: true,
			errSign: "unknown field",
		},
		{
			name: "Failure - Request body size exceeds maxBytes limit",
			args: args{
				body:     io.NopCloser(strings.NewReader(`{"name":"John","age":30}`)),
				maxBytes: 10, // Límite ultra bajo para forzar el truncamiento
			},
			want:    DummyStruct{},
			wantErr: true,
			// Al truncar el stream, el JSON queda incompleto y el decoder nativo lanza EOF o error sintáctico
			errSign: "unexpected EOF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Invocamos la función genérica explicitando el tipo DummyStruct
			got, err := Decode[DummyStruct](tt.args.body, tt.args.maxBytes)

			// 1. Validar si esperábamos un error o no
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			// 2. Si hay error, verificar que el mensaje sea el correcto (Error Signature)
			if tt.wantErr && tt.errSign != "" {
				if !strings.Contains(err.Error(), tt.errSign) {
					t.Errorf("Decode() error message = %q, expected to contain %q", err.Error(), tt.errSign)
				}
			}

			// 3. Si no hay error, validar que el objeto mapeado sea idéntico al esperado
			if !tt.wantErr {
				if got.Name != tt.want.Name || got.Age != tt.want.Age {
					t.Errorf("Decode() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEncodeNoContent(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "Success - Emits HTTP 201 Created status with no body",
			statusCode: http.StatusCreated,
		},
		{
			name:       "Success - Emits HTTP 204 No Content status with no body",
			statusCode: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			EncodeNoContent(rec, tt.statusCode)

			if rec.Code != tt.statusCode {
				t.Errorf("EncodeNoContent() status = %d, want %d", rec.Code, tt.statusCode)
			}

			if rec.Body.Len() != 0 {
				t.Errorf("EncodeNoContent() wrote an unexpected body stream: %q", rec.Body.String())
			}
		})
	}
}
func TestEncodeError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		inputErr   error
		wantInBody string
	}{
		{
			name:       "Success - Standard error serialized into message field",
			statusCode: http.StatusInternalServerError,
			inputErr:   errors.New("internal database error"),
			wantInBody: `"message":"internal database error"`,
		},
		{
			name:       "Success - ValidationError string representation serialized into message field",
			statusCode: http.StatusUnprocessableEntity,
			inputErr: ValidationError{
				Fields: map[string]string{
					"center":  "required",
					"vehicle": "must be one of: sedan, suv",
				},
			},
			wantInBody: "validation failed:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			EncodeError(rec, tt.statusCode, tt.inputErr)

			if rec.Code != tt.statusCode {
				t.Errorf("EncodeError() code = %d, want %d", rec.Code, tt.statusCode)
			}

			gotContentType := rec.Header().Get("Content-Type")
			wantContentType := "application/json; charset=utf-8"
			if gotContentType != wantContentType {
				t.Errorf("EncodeError() content-type = %q, want %q", gotContentType, wantContentType)
			}

			bodyStr := rec.Body.String()
			if !strings.Contains(bodyStr, tt.wantInBody) {
				t.Errorf("EncodeError() body = %q, expected to contain %q", bodyStr, tt.wantInBody)
			}

			var resp ErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				t.Fatalf("EncodeError() produced an invalid JSON payload: %v", err)
			}
		})
	}
}
