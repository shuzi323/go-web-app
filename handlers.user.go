package main

import (
	"github.com/gin-gonic/gin"
	//"math/rand"
	"net/http"
	//"strconv"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
)

// 用于token， 仿照jwt
type Payload struct {
	IssueTime  string `json: "issue_time"` //发行时间
	Expiration int    `json:"expiration_time"`
	Username   string `json: "username"`
	UserID     int    `json: "user_id"`
}

const (
	Salt = "secret"
)

//登录
func showLoginPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Login"}, "login.html")
}

func performLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	OK, userID := isUserValid(username, password)
	if OK {
		token, err := createToken(username, userID, Salt)
		if err != nil {
			log.Println(err)
		}
		c.SetCookie("token", token, 3600, "", "", false, true)
		c.Set("is_logged_in", true)
		storeToken(token, userID)
		render(c, gin.H{
			"title": "Successful Login"}, "login-successful.html")
	} else {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentialsprovided"})
	}
}

func logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Set("is_logged_in", false)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

//注册
func showRegistrationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Register"}, "register.html")
}

func register(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if user, err := registerNewUser(username, password); err == nil {
		token, err := createToken(user.Username, user.ID, Salt)
		if err != nil {
			log.Println(err)
		}
		c.SetCookie("token", token, 3600, "", "", false, true)
		storeToken(token, user.ID)
		c.Set("is_logged_in", true)
		render(c, gin.H{
			"title": "Successful registration & Login"}, "login-successful.html")

	} else {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"ErrorTitle":   "Registration Failed",
			"ErrorMessage": err.Error()})

	}
}

// 产生token
func createToken(username string, userID int, salt string) (string, error) {
	header := `{"typ": "JWT", "alg": "HS256"}`
	header_64 := base64.StdEncoding.EncodeToString([]byte(header))
	payload := Payload{IssueTime: time.Now().Format("2006-01-02 15:04:05"), Expiration: 3600, Username: username, UserID: userID}
	user_json, err := json.Marshal(payload) //payload 转为json []byte类型
	if err != nil {
		return "", err
	}
	user_64 := base64.StdEncoding.EncodeToString(user_json)
	header_user_64 := header_64 + "." + user_64
	HS256 := HmacSha256([]byte(header_user_64), []byte(Salt)) //加密后的
	signature_64 := base64.StdEncoding.EncodeToString(HS256)
	return header_user_64 + "." + signature_64, nil
}

// 获取token中的用户ID
func getTokenUserID(token string) (int, error) {

	user_64 := strings.Split(token, ".")
	if len(user_64) != 3 {
		return 0, errors.New("不合法的token")
	}
	user_json, _ := base64.StdEncoding.DecodeString(user_64[1])
	var user Payload
	json.Unmarshal(user_json, &user)
	return user.UserID, nil
}

// 加密
func HmacSha256(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	HS256 := mac.Sum(nil)
	return HS256
}

// 判断输入的message是否正确
func CheckMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}
