package users

import (

	"gogam/database"
	"gogam/database/model"
	"gogam/utility"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	Gorm "github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func getUserByID(c *gin.Context) {
		db := c.MustGet("db").(*Gorm.DB)
		stringID :=  c.Param("id")
		u64, _ := strconv.ParseUint(stringID, 10,32)
		 u := uint(u64)

	var searchUser model.User
	db.Model(&model.User{Model: gorm.Model{ID: u}}).Find(&searchUser)
	c.JSON(200,model.Envelope{Status: "success",Message: searchUser})
}
func getAllUsers(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	var searchUser []model.User
	db.Find(&searchUser)
	c.JSON(200,model.Envelope{Status: "success",Message: searchUser})
}

func getAllCharactersFromUser(c *gin.Context){
	db := c.MustGet("db").(*Gorm.DB)
	stringID :=  c.Param("id")
	u64, _ := strconv.ParseUint(stringID, 10,32)
	var searchCharacters []model.Character
	
	db.Where("user_id = ?",u64).Find(&searchCharacters)
	c.JSON(200,model.Envelope{Status: "success",Message: searchCharacters})
}

func createCharacter(c *gin.Context)  {
	db := c.MustGet("db").(*Gorm.DB)

	stringID :=  c.Param("id")
	u64, _ := strconv.ParseUint(stringID, 10,32)
	u := uint(u64)
	
	var newCharacter model.Character
	var linkedUser model.User
	c.BindJSON(&newCharacter)
	log.Infof("user id: %d and requested new Character %s" , u, utility.PrettyJSON(newCharacter))
	
	db.Model(&model.User{Model: gorm.Model{ID: u}}).Find(&linkedUser)
	linkedUser.Character = append(linkedUser.Character, newCharacter)
	err := database.UpdateEntity(&linkedUser,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,model.Envelope{Status: "failed",Error: err.Error()})
			return
		}

	c.JSON(200,model.Envelope{Status: "success",Message: newCharacter})
}


func newUser(c *gin.Context)  {
	db := c.MustGet("db").(*Gorm.DB)
	var newUser model.User
	c.BindJSON(&newUser)
	newUser.Password = utility.HashPassword(newUser.Password)
	err := database.CreateEntity(&newUser,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,model.Envelope{Status: "failed",Error: err.Error()})
			return
		}
	c.JSON(200,model.Envelope{Status: "success",Message: newUser})
}