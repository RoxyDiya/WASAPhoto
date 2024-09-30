package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// PROFILE

	rt.router.POST("/session", rt.doLogin)
	rt.router.PUT("/user/:userId/update-username", rt.authWrapper(rt.setMyUserName))
	rt.router.GET("/user/:userId/profile-page/:username", rt.authWrapper(rt.getUserProfile))
	rt.router.GET("/user/:userId/search/:username", rt.authWrapper(rt.searchUser))

	// SOCIAL ACTIONS

	rt.router.PUT("/user/:userId/follow/:username", rt.authWrapper(rt.followUser))
	rt.router.DELETE("/user/:userId/follow/:username", rt.authWrapper(rt.unfollowUser))
	rt.router.PUT("/user/:userId/ban/:username", rt.authWrapper(rt.banUser))
	rt.router.DELETE("/user/:userId/ban/:username", rt.authWrapper(rt.unbanUser))

	// PHOTOS INERACTIONS

	rt.router.GET("/user/:userId/photos/", rt.authWrapper(rt.getMyStream))
	rt.router.POST("/user/:userId/photos/", rt.authWrapper(rt.uploadPhoto))
	rt.router.GET("/user/:userId/photos/:photoId/", rt.authWrapper(rt.getPhoto))
	rt.router.DELETE("/user/:userId/photos/:photoId/", rt.authWrapper(rt.deletePhoto))
	rt.router.PUT("/user/:userId/photos/:photoId/likes/:authenticatedUserId", rt.authWrapper(rt.likePhoto))
	rt.router.DELETE("/user/:userId/photos/:photoId/likes/:authenticatedUserId", rt.authWrapper(rt.unlikePhoto))


	// COMMENTS
	rt.router.GET("/user/:userId/photos/:photoId/comments/", rt.authWrapperNoPath(rt.getPhotoComments))
	rt.router.POST("/user/:userId/photos/:photoId/comments/", rt.authWrapperNoPath(rt.commentPhoto))
	rt.router.DELETE("/user/:userId/photos/:photoId/comments/:commentId", rt.authWrapperNoPath(rt.deleteComment))

	return rt.router
}