package handler

import "database/sql"

// Handler holds db connection
type Handler struct {
	DB *sql.DB
}

// NewHandler returns Handler with specified db connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}
