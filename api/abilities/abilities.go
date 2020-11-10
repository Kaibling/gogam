package abilities

import (
	"gogam/middleware"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes( r *gin.RouterGroup) *gin.RouterGroup {
	abilities := r.Group("/abilities")
	{
	abilities.POST("", middleware.CheckAuthentication, newAbility)
	//users.GET("", getAllItems)
	//users.GET("/:id", getItemByID)
	abilities.GET("/:id",  middleware.CheckAuthentication, getAbilityByID)
	//users.GET("/:id/characters", getAllCharactersFromUser)
	//users.POST("/:id/characters", createCharacter)
	}
	return r
}