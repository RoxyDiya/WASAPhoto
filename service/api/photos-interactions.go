package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

// Utility to handle error response in a compact way
func handleError(w http.ResponseWriter, err error, defaultStatus int, customMessage string) bool {
	if err != nil {
		if customMessage != "" {
			ReturnCustomMessage(w, customMessage, defaultStatus)
		} else {
			ReturnInternalServerError(w, err)
		}
		return true
	}
	return false
}

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photo, err := io.ReadAll(r.Body)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo data") || len(photo) == 0 {
		return
	}

	if handleError(w, rt.db.PostPhoto(photo, token), http.StatusInternalServerError, "") {
		return
	}

	ReturnCreatedMessage(w)
}

func (rt *_router) deletePhoto(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnNotFoundError(w)
		return
	}

	if !rt.db.CheckPhotoOwner(token, photoId) {
		ReturnForbiddenMessage(w)
		return
	}

	if handleError(w, rt.db.DeletePhoto(token, photoId), http.StatusInternalServerError, "") {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getPhoto(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	owner, err := rt.db.GetPhotoOwner(photoId)
	if handleError(w, err, http.StatusInternalServerError, "") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnNotFoundError(w)
		return
	}

	if rt.db.CheckBan(owner, token) {
		ReturnForbiddenMessage(w)
		return
	}

	photo, err := rt.db.GetImage(photoId)
	if handleError(w, err, http.StatusInternalServerError, "") {
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(photo)
}

func (rt *_router) likePhoto(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	pathOwner, err := strconv.ParseInt(p.ByName("userId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid user ID") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnNotFoundError(w)
		return
	}

	if owner, _ := rt.db.GetPhotoOwner(photoId); owner != pathOwner {
		ReturnBadRequestCustomMessage(w)
		return
	}

	if rt.db.CheckBan(pathOwner, token) {
		ReturnForbiddenMessage(w)
		return
	}

	if rt.db.CheckLike(token, photoId) {
		ReturnConflictMessage(w)
		return
	}

	if handleError(w, rt.db.LikePhoto(token, photoId), http.StatusInternalServerError, "") {
		return
	}

	ReturnCreatedMessage(w)
}

func (rt *_router) unlikePhoto(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	pathOwner, err := strconv.ParseInt(p.ByName("userId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid user ID") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnNotFoundError(w)
		return
	}

	if owner, _ := rt.db.GetPhotoOwner(photoId); owner != pathOwner {
		ReturnBadRequestCustomMessage(w)
		return
	}

	if rt.db.CheckBan(pathOwner, token) {
		ReturnForbiddenMessage(w)
		return
	}

	if !rt.db.CheckLike(token, photoId) {
		ReturnConflictMessage(w)
		return
	}

	if handleError(w, rt.db.UnlikePhoto(token, photoId), http.StatusInternalServerError, "") {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnNotFoundError(w)
		return
	}

	if owner, _ := rt.db.GetPhotoOwner(photoId); rt.db.CheckBan(owner, token) {
		ReturnForbiddenMessage(w)
		return
	}

	var comment Comment
	if handleError(w, json.NewDecoder(r.Body).Decode(&comment), http.StatusBadRequest, "Invalid comment data") {
		return
	}

	newId, err := rt.db.CommentPhoto(token, photoId, comment.Comment)
	if handleError(w, err, http.StatusInternalServerError, "") {
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreatedCommentMessage{CommentId: newId})
}

func (rt *_router) getPhotoComments(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photoId, err := strconv.ParseInt(p.ByName("photoId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid photo ID") {
		return
	}

	if !rt.db.CheckPhotoExistence(photoId) {
		ReturnBadRequestMessage(w, "Photo not found")
		return
	}

	if owner, _ := rt.db.GetPhotoOwner(photoId); rt.db.CheckBan(owner, token) {
		ReturnForbiddenMessage(w)
		return
	}

	comments, err := rt.db.GetPhotoComments(photoId)
	if handleError(w, err, http.StatusInternalServerError, "") {
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (rt *_router) deleteComment(w http.ResponseWriter, _ *http.Request, p httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	commentId, err := strconv.ParseInt(p.ByName("commentId"), 10, 64)
	if handleError(w, err, http.StatusBadRequest, "Invalid comment ID") {
		return
	}

	if owner, _ := rt.db.GetCommentOwner(commentId); owner != token {
		ReturnForbiddenMessage(w)
		return
	}

	if handleError(w, rt.db.DeleteComment(commentId), http.StatusInternalServerError, "") {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getMyStream(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, token int64) {
	w.Header().Set("Content-Type", "application/json")

	photos, err := rt.db.GetMyStream(token)
	if handleError(w, err, http.StatusInternalServerError, "") {
		return
	}

	json.NewEncoder(w).Encode(photos)
}
