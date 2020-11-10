package characters

import (
	"gogam/middleware"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes( r *gin.RouterGroup) *gin.RouterGroup {
	characters := r.Group("/characters")
	{
	characters.PUT("/:cid/abilities/:aid",  middleware.CheckAuthentication , addAbilityToCharacter)
	characters.GET("/:cid", middleware.CheckAuthentication,getCharacterByID)
	//users.GET("/:id", getItemByID)
	//abilities.GET("/:id", getAbilityByID)
	//users.GET("/:id/characters", getAllCharactersFromUser)
	//users.POST("/:id/characters", createCharacter)
	}
	return r
}
