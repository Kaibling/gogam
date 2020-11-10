package characters

import (

	"gogam/database"
	"gogam/database/model"
	"gogam/utility"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	Gorm "github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)


func addAbilityToCharacter(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	characterID := utility.StringToUint( c.Param("id"))
	abilityID := utility.StringToUint( c.Param("id"))
	var resultAbility model.Ability
	var linkedCharacter model.Character

	db.Model(&model.Ability{Model: gorm.Model{ID: abilityID}}).Find(&resultAbility)
	db.Model(&model.Character{Model: gorm.Model{ID: characterID}}).Find(&linkedCharacter)
	linkedCharacter.Abilities = append(linkedCharacter.Abilities, resultAbility)
	err := database.UpdateEntity(&linkedCharacter,db)
	if err != nil {
		log.Warnln(err.Error())
		c.JSON(200,model.Envelope{Status: "failed",Error: err.Error()})
		return
	}

	c.JSON(200,model.Envelope{Status: "success",Message: resultAbility})
}


func getCharacterByID(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	characterID := utility.StringToUint( c.Param("cid"))
	var resultCharacter model.Character

	db.Model(&model.Character{Model: gorm.Model{ID: characterID}}).Preload("Abilities").Find(&resultCharacter)
	c.JSON(200,model.Envelope{Status: "success",Message: resultCharacter})
}