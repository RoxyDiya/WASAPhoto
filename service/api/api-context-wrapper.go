package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, int64)

func (rt *_router) authWrapper(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("content-type", "application/json")

		// Extract the token from the Authorization header
		token, err := ExtractToken(r)
		if err != nil || token == -1 {
			w.WriteHeader(http.StatusUnauthorized)
			rt.baseLogger.Errorf("No Token: %v", err)
			res := Message{
				Message: "No Token in the Header",
			}
			err = json.NewEncoder(w).Encode(res)
			ReturnInternalServerError(w, err)
			return
		}

		// Check if the token is valid
		if !rt.db.CheckToken(token) {
			w.WriteHeader(http.StatusNotFound)
			rt.baseLogger.Errorf("Not Active Token: %v", err)
			res := Message{
				Message: "Not Active Token",
			}
			err = json.NewEncoder(w).Encode(res)
			ReturnInternalServerError(w, err)
			return
		}

		// Prepare to handle path parameters (if they exist)
		pathParameters := [3]string{"", "", ""}
		for i, param := range ps {
			pathParameters[i] = param.Key
		}

		// Handle token in the path parameters (e.g., userId or authenticatedUserId)
		var pathToken int64 = token // Default to token from Authorization header
		if pathParameters[0] == "userId" {
			if pathParameters[2] != "authenticatedUserId" {
				pathToken = ExtractTokenFromPath(w, ps, "userId")
			} else if pathParameters[2] == "authenticatedUserId" {
				pathToken = ExtractTokenFromPath(w, ps, "authenticatedUserId")
			}
		}

		// Compare path token and Authorization header token if applicable
		if pathToken != token {
			w.WriteHeader(http.StatusForbidden)
			rt.baseLogger.Errorf("The path tokens and the auth token aren't equal: %v", err)
			res := Message{
				Message: "The path tokens and the auth token aren't equal",
			}
			err = json.NewEncoder(w).Encode(res)
			ReturnInternalServerError(w, err)
			return
		}

		// Call the next handler function with the verified token
		fn(w, r, ps, pathToken)
	}
}
