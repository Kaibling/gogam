package api

import ("github.com/gin-gonic/gin"
		"gogam/api/users"

)

func ApplyRoutes( r *gin.Engine) *gin.Engine {
	userGroup := r.Group("/api")
	{
		users.ApplyRoutes(userGroup)
	}
	return r
}