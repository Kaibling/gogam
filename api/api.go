package api

import (
	"gogam/api/abilities"
	"gogam/api/authentication"
	"gogam/api/characters"
	"gogam/api/users"

	"github.com/gin-gonic/gin"
)

func ApplyRoutes( r *gin.Engine) *gin.Engine {
	api := r.Group("/api")
	{
		users.ApplyRoutes(api)
		characters.ApplyRoutes(api)
		abilities.ApplyRoutes(api)
		authentication.ApplyRoutes(api)
	}
	return r
}