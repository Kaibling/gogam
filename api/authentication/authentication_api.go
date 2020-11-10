package authentication

import (
	"gogam/database/model"
	"gogam/utility"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	Gorm "github.com/jinzhu/gorm"
)

type userLogin struct{
	Name string
	Password string
}

func login(c *gin.Context) {
	db := c.MustGet("db").(*Gorm.DB)
	hmacSampleSecret := c.MustGet("hmacSecret").([]byte)
	var user userLogin
	var savedUser model.User

	c.BindJSON(&user)
	
	db.Model(&model.User{Name: user.Name}).Find(&savedUser)

	if savedUser.Password != utility.HashPassword(user.Password) {
		c.JSON(401,model.Envelope{Status: "failed",Error: "Access Denied"})
		c.Abort()
    return
	}	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": user.Name,
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		c.JSON(500,model.Envelope{Status: "failed",Error: err.Error()})
		c.Abort()
    return
	}
 
	c.JSON(200,model.Envelope{Status: "success",Message: tokenString})

}

func check(c *gin.Context) {

}