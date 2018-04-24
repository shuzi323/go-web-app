package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

var router *gin.Engine
var db *sql.DB

//图标
var icon []byte

func main() {
	var err error
	//数据库
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/restful?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	//加载图标
	icon, err = ioutil.ReadFile("static/img/favicon.ico")
	if err != nil {
		log.Println(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	// 加载HTML文件
	router.LoadHTMLGlob("templates/*")
	initializeRoutes()
	router.Run()
}

func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		c.XML(http.StatusOK, data["payload"])
	default:
		c.HTML(http.StatusOK, templateName, data)
	}
}
