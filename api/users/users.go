package users
import "github.com/gin-gonic/gin"

func ApplyRoutes( r *gin.RouterGroup) *gin.RouterGroup {
	users := r.Group("/users")
	{
	users.POST("", newUser)
	users.GET("", getAllUsers)
	users.GET("/:id", getUserByID)
	users.GET("/:id/characters", getAllCharactersFromUser)
	users.POST("/:id/characters", createCharacter)
	}
	return r
}