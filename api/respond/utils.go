// HSE estimation API.
//
// Lessons estimation project
//
//     Schemes: http
//     Host: localhost
//     BasePath: /api/v1/
//     Version: 0.0.1
//     Contact: Alexandr Komarov <amkomarov@edu.hse.ru>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package respond

import (
	"context"
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

func WithJSON(ctx context.Context, w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		logger.WithError(err).Error("write answer")
	}
}

func WithErrWrapJSON(ctx context.Context, w http.ResponseWriter, code int, in error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	logger.WithError(in).Error(in.Error())
	err := json.NewEncoder(w).Encode(in.Error())
	if err != nil {
		logger.WithError(err).Error("write answer")
	}
}

// Error while processing request
//
// swagger:response errorResponse
type errorResponse struct {
	// In: body
	Body struct {
		// Error message
		//
		Error string `json:"error"`
	}
}

// File not found
//
// swagger:response notFoundResponse
type notFoundResponse struct {
}

// Empty response
//
// swagger:response statusOkResponse
type statusOkResponse struct {
	// in: body
	Body struct {
	}
}

// Bad request
//
// swagger:response badRequestResponse
type badRequestResponse struct {
	// In: body
	Body struct {
		// Error message
		//
		Error string `json:"error"`
	}
}