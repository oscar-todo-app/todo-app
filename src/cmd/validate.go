package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func validateIDParam(w http.ResponseWriter, idParam string) (int, error) {
	// Validate if idParam is a non-empty string
	if idParam == "" {
		handleError(
			w,
			"Invalid 'id' parameter",
			"invalid or missing 'id' parameter",
			http.StatusBadRequest,
		)
		return 0, fmt.Errorf("invalid or missing 'id' parameter")
	}

	// Validate if idParam is a valid integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		handleError(
			w,
			"Invalid 'id' parameter",
			"invalid or missing 'id' parameter",
			http.StatusBadRequest,
		)
		return 0, fmt.Errorf("invalid or missing 'id' parameter")
	}

	// Validate if id is greater than 0 (assuming your IDs are positive integers)
	if id <= 0 {
		handleError(w, "Invalid 'id' parameter", "invalid 'id' parameter", http.StatusBadRequest)
		return 0, fmt.Errorf("invalid 'id' parameter")
	}

	return id, nil
}

func handleError(w http.ResponseWriter, logMsg string, errMsg string, statusCode int) {
	log.Printf("%s: %s", logMsg, errMsg)
	http.Error(w, errMsg, statusCode)
}
