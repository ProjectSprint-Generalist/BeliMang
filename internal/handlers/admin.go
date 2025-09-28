package handlers

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
)

// AdminHandler wires admin endpoints to sqlc-generated queries.
type AdminHandler struct {
	queries *db.Queries
}

func NewAdminHandler(queries *db.Queries) *AdminHandler {
	return &AdminHandler{queries: queries}
}
