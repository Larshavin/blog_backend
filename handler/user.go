package handler

import (
	db "blog/database"
	"blog/ent"
	"blog/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var frontData ent.User
		err := c.BindJSON(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// hash password
		frontData.HashedPassword, err = user.HashPassword(frontData.HashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// user role
		frontData.Role = "visitor"

		err = db.SaveUser(&frontData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func PostLogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		type dataType struct {
			AccessToken string `json:"accessToken"`
		}
		var json dataType

		// decode token and get email
		err := c.BindJSON(&json)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = db.DeleteTokenByAccessToken(json.AccessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "logout success"})
	}
}

func PostloginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginData user.LoginData
		err := c.BindJSON(&loginData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		userInfo, err := user.CheckPassword(loginData.Email, loginData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		accessToken, refreshToken, err := user.GenerateToken(*userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// save tokens to db
		err = db.SaveToken(userInfo.Email, accessToken, refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// marshal data
		data := gin.H{"accessToken": accessToken, "refreshToken": refreshToken}

		c.JSON(http.StatusOK, data)
	}
}
