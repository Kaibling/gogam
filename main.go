package main

import (
	"encoding/json"
	"strconv"

	"github.com/kaibling/gogam/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	Gorm "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
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


func initDB() *Gorm.DB {
	
	//db, err := gorm.Open("sqlite3", ":memory:")
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Character{})
	db.AutoMigrate(&model.Ability{})
	return db
}

func init() {
db := initDB()
	defer db.Close()

	/*

		cu := &User{Name: "dragon",Password: "asd"}
		err := createEntity(cu,db)
		if err != nil {
			fmt.Println(err.Error())
		}
		db.Create(&User{Name: "dragon"})
		cu.Password = "supergeil"
		err = updateEntity(cu,db)
				if err != nil {
			fmt.Println(err.Error())
		}
		
		var derUser User
		db.Model(&User{Name: "dragon"}).Find(&derUser)
		prettyJSON(derUser)
*/
}

func getUserByID(c *gin.Context) {
		db := c.MustGet("db").(*Gorm.DB)
		stringID :=  c.Param("id")
		u64, _ := strconv.ParseUint(stringID, 10,32)
		 u := uint(u64)

	var searchUser model.User
	db.Model(&model.User{Model: gorm.Model{ID: u}}).Find(&searchUser)

	//stringUser,_ := json.Marshal(searchUser)
	c.JSON(200,envelope{Status: "success",Message: searchUser})
}
func getAllUsers(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	var searchUser []model.User
	db.Find(&searchUser)

	//stringUser,_ := json.Marshal(searchUser)
	c.JSON(200,envelope{Status: "success",Message: searchUser})
}

func CreateCharacter(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	var newCharacter model.Character
	c.BindJSON(&newCharacter)
	if newCharacter.UserID == 0 {
		log.Warnln("no UserID")
			c.JSON(200,envelope{Status: "failed",Error: "no UserID"})
			return
	}
	err := createEntity(&newCharacter,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,envelope{Status: "failed",Error: err.Error()})
			return
		}

	//stringUser,_ := json.Marshal(newCharacter)
	c.JSON(200,envelope{Status: "success",Message: newCharacter})

}

func getAllCharactersFromUser(c *gin.Context){
	db := c.MustGet("db").(*Gorm.DB)

	stringID :=  c.Param("id")
	u64, _ := strconv.ParseUint(stringID, 10,32)
	u := uint(u64)


	var searchCharacters []model.Character
	db.Debug().Model(&model.User{Model: gorm.Model{ID: u}}).Association("characters").Find(&searchCharacters)

	//stringUser,_ := json.Marshal(searchCharacters)
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

	db.Debug().Model(&model.User{Model: gorm.Model{ID: u}}).Find(&linkedUser)
	linkedUser.Character = append(linkedUser.Character, newCharacter)
	err := updateEntity(&linkedUser,db)
	if err != nil {
			log.Warnln(err.Error())
			c.JSON(200,envelope{Status: "failed",Error: err.Error()})
			return
		}

	stringUser,_ := json.Marshal(newCharacter)
	c.JSON(200,envelope{Status: "success",Message: stringUser})
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

	//stringUser,_ := json.Marshal(newUser)
	c.JSON(200,envelope{Status: "success",Message: newUser})
}


func main() {
	db := initDB()
	defer db.Close()
	
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context){
		c.Set("db",db)
	})
	r.POST("/users", newUser)
	r.GET("/users", getAllUsers)
	r.GET("/users/:id", getUserByID)
	r.GET("/users/:id/characters", getAllCharactersFromUser)
	r.POST("/users/:id/characters", createCharacter)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() 
}

func prettyJSON(object interface{}) {
	a, _ := json.MarshalIndent(object, "", " ")
	log.Infoln(string(a))
}
