package middleware

import (
	"gogam/database/model"
	//"gogam/utility"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
	"fmt"
)
//CheckAuthentication handles authentication in the middleware
func CheckAuthentication(c *gin.Context) {

	hmacSampleSecret := c.MustGet("hmacSecret").([]byte)
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.JSON(403,model.Envelope{Status: "failed",Error: "Could not find Authorization header"})
			c.Abort()
			return

	}
	tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			c.JSON(403,model.Envelope{Status: "failed",Error:"Could not find bearer token in Authorization header"})
			c.Abort()
			return
		}

token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }
    return hmacSampleSecret, nil
})
if err != nil {

	log.Infof("Bad Thing happened: %s",err.Error())
	c.JSON(500,model.Envelope{Status: "failed",Error:err.Error()})
	c.Abort()
	return

}

claims, ok := token.Claims.(jwt.MapClaims)
if !ok || !token.Valid {
	log.Infoln("token invalid")
	c.JSON(403,model.Envelope{Status: "failed",Error:"token invalid"})
	c.Abort()
	return
}
/*
if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

}
*/
	c.Set("userName",claims["name"])
	c.Next()
}