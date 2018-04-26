package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if !loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func ensureNotLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		loggedIn := loggedInInterface.(bool)
		if loggedIn {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func setUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" || checkToken(token) {
			c.Set("is_logged_in", true)
		} else {
			c.Set("is_logged_in", false)
		}
	}
}

// 验证用户token和数据库保存的是否相同
func checkToken(token string) bool {
	userID, err := getTokenUserID(token)
	if err != nil {
		return false
	}
	var databaseToken string
	err = db.QueryRow("SELECT token_data FROM token WHERE user_id = ?", userID).Scan(&databaseToken)
	if err != nil {
		log.Println(err) //数据库不存在该ID
		return false
	} else if databaseToken == token {
		return true
	}
	return false
}
