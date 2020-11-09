package users
import (
	"gogam/database/model"
	//"gogam/model"

	"github.com/gin-gonic/gin"
	Gorm "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm"
		"strconv"
			log "github.com/sirupsen/logrus"
				"encoding/json"
)

type envelope struct {
	Status string `json:status`
	Error string `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

func createEntity(model interface{},db *Gorm.DB) error{
	result := db.Create(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func updateEntity(model interface{},db *Gorm.DB) error{
	result := db.Save(model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}


func prettyJSON(object interface{}) string {
	a, _ := json.MarshalIndent(object, "", " ")
	return string(a)
}







func getUserByID(c *gin.Context) {
		db := c.MustGet("db").(*Gorm.DB)
		stringID :=  c.Param("id")
		u64, _ := strconv.ParseUint(stringID, 10,32)
		 u := uint(u64)

	var searchUser model.User
	db.Model(&model.User{Model: gorm.Model{ID: u}}).Find(&searchUser)
	c.JSON(200,envelope{Status: "success",Message: searchUser})
}
func getAllUsers(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	var searchUser []model.User
	db.Find(&searchUser)
	c.JSON(200,envelope{Status: "success",Message: searchUser})
}

func getAllCharactersFromUser(c *gin.Context){
	db := c.MustGet("db").(*Gorm.DB)
	stringID :=  c.Param("id")
	u64, _ := strconv.ParseUint(stringID, 10,32)
	var searchCharacters []model.Character
	
	db.Where("user_id = ?",u64).Find(&searchCharacters)
	c.JSON(200,envelope{Status: "success",Message: searchCharacters})


}

func createCharacter(c *gin.Context)  {
	db := c.MustGet("db").(*Gorm.DB)

	stringID :=  c.Param("id")
	u64, _ := strconv.ParseUint(stringID, 10,32)
	u := uint(u64)
	
	var newCharacter model.Character
	var linkedUser model.User
	c.BindJSON(&newCharacter)
	log.Infof("user id: %d and requested new Character %s" , u,prettyJSON(newCharacter))
	
	db.Model(&model.User{Model: gorm.Model{ID: u}}).Find(&linkedUser)
	linkedUser.Character = append(linkedUser.Character, newCharacter)
	err := updateEntity(&linkedUser,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,envelope{Status: "failed",Error: err.Error()})
			return
		}

	c.JSON(200,envelope{Status: "success",Message: newCharacter})
}


func newUser(c *gin.Context)  {
	db := c.MustGet("db").(*Gorm.DB)
	var newUser model.User
	c.BindJSON(&newUser)
	err := createEntity(&newUser,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,envelope{Status: "failed",Error: err.Error()})
			return
		}
	c.JSON(200,envelope{Status: "success",Message: newUser})
}