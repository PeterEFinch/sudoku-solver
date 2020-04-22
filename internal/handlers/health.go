package handlers

import (
	"net/http"
)

// Health handles requests for the health of the service.
func Health(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
}
