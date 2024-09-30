package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) handleUserAction(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64, action func(int64, int64) error, reverseAction bool, errMsg string) {
	w.Header().Set("content-type", "application/json")

	// Get username from the path
	username := p.ByName("username")

	// Check if the username respects the regex
	if !CheckUsernameRegex(w, username) {
		ReturnBadRequestMessage(w, nil)
		return
	}

	// Get the token of the username
	token2, err := rt.db.GetUserTokenOnly(username)
	if err != nil {
		ReturnNotFoundError(w)
		return
	}

	// Check if the user is trying to act on himself
	if token == token2 {
		ReturnForbiddenMessage(w)
		return
	}

	// Check if the user is banned from the user (applicable in follow/unfollow scenarios)
	if ban, err := rt.db.CheckBan(token2, token); err != nil || ban {
		ReturnForbiddenMessage(w)
		return
	}

	// Check if the action is already performed or not (e.g., already followed/unfollowed)
	check, err := rt.db.CheckFollow(token, token2)
	if err != nil {
		ReturnInternalServerError(w, err)
		return
	}

	if check == reverseAction {
		// Return appropriate message based on the context (conflict if already exists, forbidden otherwise)
		if reverseAction {
			ReturnForbiddenMessage(w)
		} else {
			ReturnConflictMessage(w)
		}
		return
	}

	// Perform the desired action (follow, unfollow, ban, unban)
	if err := action(token, token2); err != nil {
		ReturnInternalServerError(w, err)
		return
	}

	// Set the appropriate response based on action type
	if reverseAction {
		w.WriteHeader(http.StatusNoContent)
	} else {
		ReturnCreatedMessage(w)
	}
}

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64) {
	rt.handleUserAction(w, r, p, token, rt.db.AddFollow, false, "Already following")
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64) {
	rt.handleUserAction(w, r, p, token, rt.db.RemoveFollow, true, "Not following")
}

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64) {
	rt.handleUserAction(w, r, p, token, rt.db.AddBan, false, "Already banned")
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64) {
	rt.handleUserAction(w, r, p, token, rt.db.RemoveBan, true, "Not banned")
}
