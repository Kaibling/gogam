package abilities

import (
	"strconv"
	"gogam/database/model"
	"gogam/database"
	 "github.com/gin-gonic/gin"
	Gorm "github.com/jinzhu/gorm"
	 "github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)


func newAbility(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	var newAbility model.Ability
	c.BindJSON(&newAbility)
	err := database.CreateEntity(&newAbility,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,model.Envelope{Status: "failed",Error: err.Error()})
			return
		}
	c.JSON(200,model.Envelope{Status: "success",Message: newAbility})
}

func getAllItems(c *gin.Context) {
	
}


func getAbilityByID(c *gin.Context) {

			db := c.MustGet("db").(*Gorm.DB)
		stringID :=  c.Param("id")
		u64, _ := strconv.ParseUint(stringID, 10,32)
		 u := uint(u64)

	var searchAbility model.Ability
	db.Model(&model.Ability{Model: gorm.Model{ID: u}}).Find(&searchAbility)
	c.JSON(200,model.Envelope{Status: "success",Message: searchAbility})
	
}