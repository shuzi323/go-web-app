package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// 首页
func showIndexPage(c *gin.Context) {
	var page Page
	var err error
	page.PageNow, err = strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	page.PageSize = 8
	page.myPage()
	if page.PageNow > page.TotalPage || page.PageNow <= 0 { //页面不存在
		c.AbortWithStatus(http.StatusNotFound)
	}
	//articles, err := getAllArticles()  分页后用不到全部获取
	var articles []Article
	articles, err = page.getArticles()
	if err != nil {
		log.Fatalln(err)
	}
	render(c, gin.H{
		"title":   "Home Page",
		"page":    page,
		"payload": articles}, "index.html")
}

//获取文章
func getArticle(c *gin.Context) {
	if articleID, err := strconv.Atoi(c.Param("article_id")); err == nil {
		if article, err := getArticleByID(articleID); err == nil {
			render(c, gin.H{
				"title":   article.Title,
				"payload": article}, "article.html")
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// 打开新建文章页面
func showArticleCreationPage(c *gin.Context) {
	render(c, gin.H{
		"title": "Create New Article"}, "create-article.html")
}

//新建文章
func createArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	//获取用户ID
	token, _ := c.Cookie("token")
	userID, _ := getTokenUserID(token)
	if a, err := createNewArticle(title, content, userID); err == nil {
		render(c, gin.H{
			"title":   "Submission Successful",
			"payload": a}, "submission-successful.html")
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

//图标
func favicon(c *gin.Context) {
	c.Data(http.StatusOK, "image/x-icon", icon)
}
