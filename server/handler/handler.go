package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Handler holds db connection
type Handler struct {
	DB *sql.DB
}

// NewHandler returns Handler with specified db connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

const (
	// MaxRequestBodySize limits the size of request bodies to prevent DoS attacks
	MaxRequestBodySize = 1024 * 1024 // 1MB
)

// parseRequestBody reads and parses JSON request body
func parseRequestBody(r *http.Request, v interface{}) error {
	reqlen, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return fmt.Errorf("cannot get Content-Length: %w", err)
	}
	if reqlen > MaxRequestBodySize {
		return fmt.Errorf("request body too large: %d bytes (max %d)", reqlen, MaxRequestBodySize)
	}
	if reqlen < 0 {
		return fmt.Errorf("invalid Content-Length: %d", reqlen)
	}
	body := make([]byte, reqlen)
	n, err := r.Body.Read(body)
	if err != nil {
		if err != io.EOF || n == 0 {
			return fmt.Errorf("failed to read request body: %w", err)
		}
	}
	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return nil
}

// handleRequestError logs the error and writes an error response
func handleRequestError(w http.ResponseWriter, r *http.Request, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, "%s: %v", message, err)
	log.Printf("%v %v: %s: %v", r.Method, r.URL.Path, message, err)
}

// validateSID validates student ID format (basic validation)
func validateSID(sid string) error {
	if sid == "" {
		return fmt.Errorf("SID cannot be empty")
	}
	if len(sid) > 100 {
		return fmt.Errorf("SID too long (max 100 characters)")
	}
	return nil
}

// validateCardID validates card ID format
func validateCardID(cardID string) error {
	if cardID == "" {
		return fmt.Errorf("CardID cannot be empty")
	}
	if len(cardID) > 100 {
		return fmt.Errorf("CardID too long (max 100 characters)")
	}
	return nil
}
