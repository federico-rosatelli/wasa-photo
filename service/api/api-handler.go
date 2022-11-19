package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.Welcome)
	rt.router.POST("/", rt.AddSeen)
	rt.router.GET("/getinfo", rt.GetInfo)
	rt.router.GET("/login", rt.Login)
	rt.router.POST("/signin", rt.SignIn)
	rt.router.GET("/addphoto", rt.AddPhotoProfileGet)
	rt.router.POST("/addphoto", rt.AddPhotoProfile)
	rt.router.DELETE("/addphoto", rt.DeletePhotoProfile)
	rt.router.GET("/profile", rt.ProfileInfo)
	rt.router.POST("/profile", rt.UpdateProfileInfo)
	rt.router.GET("/profile/{id}", rt.GetBasicProfile)
	rt.router.POST("/profile/{id}", rt.AddFollowerProfile)
	rt.router.DELETE("/profile/{id}", rt.UnFollowerProfile)
	rt.router.POST("/profile/{id}/ban", rt.BanFollowerProfile)
	rt.router.DELETE("/profile/{id}/ban", rt.UnBanFollowerProfile)
	rt.router.GET("/profile/{id}/ultra", rt.GetUltraBasicProfile)
	rt.router.GET("/profile/{id}/{imageId}", rt.GetProfileImageInfo)
	rt.router.POST("/profile/{id}/{imageId}/like", rt.AddLikeProfile)
	rt.router.DELETE("/profile/{id}/{imageId}/like", rt.DeleteLikeProfile)
	rt.router.POST("/profile/{id}/{imageId}/comment", rt.AddCommentProfile)
	rt.router.DELETE("/profile/{id}/{imageId}/comment", rt.DeleteCommentProfile)
	rt.router.GET("/profile/{id}/followers", rt.GetProfileFollowers)
	rt.router.GET("/profile/{id}/followings", rt.GetProfileFollowings)
	rt.router.GET("/search", rt.SearchProfile)

	return rt.router
}
