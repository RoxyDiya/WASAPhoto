package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"structs"
)

// Helper to send a JSON response with a status code
func sendJSONResponse(w http.ResponseWriter, statusCode int, message string) error {
	w.WriteHeader(statusCode)
	res := structs.Message{Message: message}
	return json.NewEncoder(w).Encode(res)
}

// Error Messages

// ReturnInternalServerError sends a 500 Internal Server Error with an optional error message
func ReturnInternalServerError(w http.ResponseWriter, err error) {
	if err != nil {
		_ = sendJSONResponse(w, http.StatusInternalServerError, "Internal Server Error: "+err.Error())
	}
}

// ReturnNotFoundError sends a 404 Not Found response
func ReturnNotFoundError(w http.ResponseWriter) {
	err := sendJSONResponse(w, http.StatusNotFound, "Resource Not Found")
	ReturnInternalServerError(w, err)
}

// ReturnCreatedMessage sends a 201 Created response
func ReturnCreatedMessage(w http.ResponseWriter) {
	err := sendJSONResponse(w, http.StatusCreated, "Created Successfully")
	ReturnInternalServerError(w, err)
}

// ReturnBadRequestMessage sends a 400 Bad Request response
func ReturnBadRequestMessage(w http.ResponseWriter, err error) {
	if err != nil {
		_ = sendJSONResponse(w, http.StatusBadRequest, "Bad Request: "+err.Error())
		ReturnInternalServerError(w, err)
	}
}

// ReturnBadRequestCustomMessage sends a custom 400 Bad Request response
func ReturnBadRequestCustomMessage(w http.ResponseWriter) {
	err := sendJSONResponse(w, http.StatusBadRequest, "Bad Request: The request is not valid")
	ReturnInternalServerError(w, err)
}

// ReturnForbiddenMessage sends a 403 Forbidden response
func ReturnForbiddenMessage(w http.ResponseWriter) {
	err := sendJSONResponse(w, http.StatusForbidden, "Forbidden Action")
	ReturnInternalServerError(w, err)
}

// ReturnConflictMessage sends a 409 Conflict response
func ReturnConflictMessage(w http.ResponseWriter) {
	err := sendJSONResponse(w, http.StatusConflict, "This resource is already in the database")
	ReturnInternalServerError(w, err)
}

// Token Functions

// ExtractToken extracts and parses a Bearer token from the request header
func ExtractToken(r *http.Request) (int64, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return -1, errors.New("no token found")
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return -1, errors.New("invalid token format")
	}

	token, err := strconv.ParseInt(splitToken[1], 10, 64)
	if err != nil {
		return -1, errors.New("invalid token value")
	}

	return token, nil
}

// ExtractTokenFromPath extracts a token from the URL path parameters
func ExtractTokenFromPath(w http.ResponseWriter, ps httprouter.Params, paramName string) int64 {
	pathToken, err := strconv.ParseInt(ps.ByName(paramName), 10, 64)
	if err != nil {
		_ = sendJSONResponse(w, http.StatusUnauthorized, "Invalid token in the path")
		ReturnInternalServerError(w, err)
		return -1
	}
	return pathToken
}

// Validation Functions

// CheckUsernameRegex checks if the username follows a valid regex pattern
func CheckUsernameRegex(w http.ResponseWriter, username string) bool {
	match, err := regexp.MatchString(`^[a-zA-Z0-9_-]{3,16}$`, username)
	if err != nil || !match {
		_ = sendJSONResponse(w, http.StatusBadRequest, "Error matching Username regex")
		ReturnInternalServerError(w, err)
		return false
	}
	return true
}
