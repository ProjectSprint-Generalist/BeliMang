package shared

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
)

func ParseDBResult(err error) (int, string) {
	log.Printf("ParseDBResult called with error: %v (type: %T)", err, err)

	if err == nil {
		return http.StatusOK, ""
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // duplicate key error
			if strings.Contains(pgErr.ConstraintName, "username") {
				return http.StatusConflict, "Username already exists"
			}
			if strings.Contains(pgErr.ConstraintName, "email") {
				return http.StatusConflict, "Email already exists"
			}
			return http.StatusConflict, "Resource already exists"

		case "23503": // foreign key error
			return http.StatusBadRequest, "Referenced resource doesn't exist"

		case "23502": // not null error
			return http.StatusBadRequest, "Required field is missing"

		case "23414": // check error
			return http.StatusBadRequest, "Invalid data provided"

		case "42P01": // Undefined table
			return http.StatusInternalServerError, "Database schema error"

		case "08003", "08006": // Connection error
			return http.StatusServiceUnavailable, "Database temporarily unavailable"
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusRequestTimeout, "Database operation timed out"
	}

	if errors.Is(err, context.Canceled) {
		return http.StatusRequestTimeout, "Database operation was cancelled"
	}

	return http.StatusInternalServerError, "Database operation failed"
}
