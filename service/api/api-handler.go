package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.Welcome)
	rt.router.POST("/", rt.AddSeen)
	rt.router.POST("/signin", rt.SignIn)
	rt.router.POST("/addphoto", rt.AddPhotoProfile)
	rt.router.GET("/profile", rt.ProfileInfo)
	rt.router.POST("/profile", rt.UpdateProfileInfo)
	rt.router.GET("/profile/:id", rt.GetBasicProfile)
	rt.router.PUT("/profile/:id", rt.AddFollowerProfile)
	rt.router.DELETE("/profile/:id", rt.UnFollowerProfile)
	rt.router.POST("/profile/:id/ban", rt.BanFollowerProfile)
	rt.router.DELETE("/profile/:id/ban", rt.UnBanFollowerProfile)
	rt.router.GET("/profile/:id/ultra", rt.GetUltraBasicProfile)
	rt.router.GET("/profile/:id/image/:imageid", rt.GetProfileImageInfo)
	rt.router.DELETE("/profile/:id/image/:imageid", rt.DeletePhotoProfile)
	rt.router.PUT("/profile/:id/like/:imageid", rt.AddLikeProfile)
	rt.router.DELETE("/profile/:id/like/:imageId", rt.DeleteLikeProfile)
	rt.router.POST("/profile/:id/comment/:imageId", rt.AddCommentProfile)
	rt.router.DELETE("/profile/:id/comment/:imageId", rt.DeleteCommentProfile)
	rt.router.GET("/profile/:id/followers", rt.GetProfileFollowers)
	rt.router.GET("/profile/:id/followings", rt.GetProfileFollowings)
	rt.router.GET("/search", rt.SearchProfile)

	return rt.router
}
