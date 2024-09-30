package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Helper function to handle JSON response encoding and error handling.
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		ReturnInternalServerError(w, err)
	}
}

// Helper function for common error handling paths
func handleError(w http.ResponseWriter, err error, message string, statusCode int) bool {
	if err != nil {
		if message == "" {
			ReturnInternalServerError(w, err)
		} else {
			http.Error(w, message, statusCode)
		}
		return true
	}
	return false
}

// Login handler
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var username Username
	if err := json.NewDecoder(r.Body).Decode(&username); handleError(w, err, "Invalid input", http.StatusBadRequest) {
		return
	}

	if !CheckUsernameRegex(w, username.Username) {
		ReturnBadRequestMessage(w, nil)
		return
	}

	// Retrieve or create user token
	token, err := rt.db.GetUserToken(username.Username)
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}

	// Return the token
	respondWithJSON(w, http.StatusCreated, Token{Identifier: token})
}

// Set the current user's username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, _ httprouter.Params, pathToken int64) {
	var username Username
	if err := json.NewDecoder(r.Body).Decode(&username); handleError(w, err, "Invalid input", http.StatusBadRequest) {
		return
	}

	if !CheckUsernameRegex(w, username.Username) {
		ReturnBadRequestMessage(w, nil)
		return
	}

	// Check if the username already exists
	exists, err := rt.db.CheckUsernameExistence(username.Username)
	if handleError(w, err, "", http.StatusInternalServerError) || exists != 0 {
		ReturnConflictMessage(w)
		return
	}

	// Update the username
	if err := rt.db.SetUserName(pathToken, username.Username); handleError(w, err, "", http.StatusInternalServerError) {
		return
	}

	respondWithJSON(w, http.StatusOK, Message{Message: "Username updated"})
}

// Get user profile by username
func (rt *_router) getUserProfile(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	username := p.ByName("username")

	if !CheckUsernameRegex(w, username) {
		ReturnBadRequestMessage(w, nil)
		return
	}

	// Get user token
	profileToken, err := rt.db.GetUserTokenOnly(username)
	if handleError(w, err, "", http.StatusNotFound) {
		return
	}

	// Check if the user is banned
	banned, err := rt.db.CheckBan(profileToken, token)
	if handleError(w, err, "", http.StatusForbidden) || banned {
		ReturnForbiddenMessage(w)
		return
	}

	// Retrieve user profile
	profile, err := rt.db.GetUserProfile(username, token)
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}

	respondWithJSON(w, http.StatusOK, profile)
}

// Search for users by partial username
func (rt *_router) searchUser(w http.ResponseWriter, _ *http.Request, p httprouter.Params, _ int64) {
	username := p.ByName("username")

	if !CheckUsernameRegex(w, username) {
		ReturnBadRequestMessage(w, nil)
		return
	}

	// Retrieve users matching the username pattern
	users, err := rt.db.GetUsersList(username)
	if handleError(w, err, "", http.StatusInternalServerError) {
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}
