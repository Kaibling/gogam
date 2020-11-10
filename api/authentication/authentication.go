package authentication

import "github.com/gin-gonic/gin"

func ApplyRoutes( r *gin.RouterGroup) *gin.RouterGroup {
	users := r.Group("/authentication")
	{
	users.POST("/login", login)
	users.POST("check", check)
	//users.GET("/:id", getUserByID)
	//users.GET("/:id/characters", getAllCharactersFromUser)
	//users.POST("/:id/characters", createCharacter)
	}
	return r
}