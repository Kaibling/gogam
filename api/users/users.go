package users
import (
	"github.com/gin-gonic/gin"
"gogam/middleware"
)

func ApplyRoutes( r *gin.RouterGroup) *gin.RouterGroup {
	users := r.Group("/users")
	{
	users.POST("", middleware.CheckAuthentication, newUser)
	users.GET("", middleware.CheckAuthentication, getAllUsers)
	users.GET("/:id", middleware.CheckAuthentication, getUserByID)
	users.GET("/:id/characters",middleware.CheckAuthentication, getAllCharactersFromUser)
	users.POST("/:id/characters",middleware.CheckAuthentication, createCharacter)
	}
	return r
}