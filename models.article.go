package main

import (
	"errors"
	"log"
	"strings"
	"time"
)

type Article struct {
	ID       int    `json:"article_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	DateTime string `json:"datetime"`
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
}

type Page struct {
	PageNow       int   //当前页数，需初始化
	PageSize      int   //每页显示的文章条数，需初始化
	TotalCount    int   //总条数
	TotalPage     int   //总页数
	ShowPages     []int //展示的页码
	ShowFirstPage bool
	ShowLastPage  bool
}

func (p *Page) myPage() {
	err := db.QueryRow("SELECT COUNT(*) FROM article").Scan(&p.TotalCount)
	if err != nil {
		log.Println(err)
	}
	//向上取整求总页数
	if p.TotalCount%p.PageSize == 0 {
		p.TotalPage = p.TotalCount / p.PageSize
	} else {
		p.TotalPage = p.TotalCount/p.PageSize + 1
	}
	//展示的页码
	switch {
	case p.TotalPage <= 5:
		p.ShowPages = make([]int, p.TotalPage)
		for i := 0; i < p.TotalPage; i++ {
			p.ShowPages[i] = i + 1
		}
	case p.PageNow <= 2:
		p.ShowPages = make([]int, 5)
		for i := 0; i < 5; i++ {
			p.ShowPages[i] = i + 1
		}
	case p.PageNow >= p.TotalPage-1:
		p.ShowPages = make([]int, 5)
		j := 4
		for i := 0; i < 5; i++ {
			p.ShowPages[j] = p.TotalPage - i
			j--
		}
	default:
		p.ShowPages = make([]int, 5)
		for i := 0; i < 5; i++ {
			p.ShowPages[i] = p.PageNow - 2 + i
		}
	}
	//是否显示第一页和最后一页
	if p.PageNow > 3 && p.TotalPage > 5 {
		p.ShowFirstPage = true
	} else {
		p.ShowFirstPage = false
	}
	if p.PageNow < p.TotalPage-2 && p.TotalPage > 5 {
		p.ShowLastPage = true
	} else {
		p.ShowLastPage = false
	}
}
func (p *Page) haveNext() bool { //是否有下一页
	if p.PageNow == p.TotalPage {
		return false
	}
	return true
}

func (p *Page) getArticles() (articles []Article, err error) {
	articles = make([]Article, 0)
	rows, err := db.Query("SELECT article_id, title, content, datetime, username FROM article, user WHERE article.user_id = user.user_id ORDER BY article_id DESC LIMIT ?, ?", (p.PageNow-1)*p.PageSize, p.PageSize)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var perArticle Article
		rows.Scan(&perArticle.ID, &perArticle.Title, &perArticle.Content, &perArticle.DateTime, &perArticle.UserName)
		articles = append(articles, perArticle)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func getAllArticles() (articles []Article, err error) {
	articles = make([]Article, 0)
	rows, err := db.Query("SELECT article_id, title, content, datetime, username FROM article, user WHERE article.user_id = user.user_id ORDER BY article_id DESC")
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var perArticle Article
		rows.Scan(&perArticle.ID, &perArticle.Title, &perArticle.Content, &perArticle.DateTime, &perArticle.UserName)
		articles = append(articles, perArticle)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func getArticleByID(id int) (*Article, error) {
	var perArticle Article
	err := db.QueryRow("SELECT title, content, datetime, username FROM article, user WHERE article_id=? AND article.user_id = user.user_id", id).Scan(&perArticle.Title, &perArticle.Content, &perArticle.DateTime, &perArticle.UserName)
	return &perArticle, err
}

func createNewArticle(title, content string, userID int) (*Article, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("Title can'n be empty")
	}
	datetime := time.Now().Local().Format("2006-01-02 15:04:05")

	stmt, err := db.Exec("INSERT INTO article(title, content, datetime, user_id) VALUES(?, ?, ?, ?)", title, content, datetime, userID)
	if err != nil {
		return nil, err
	}
	id, err := stmt.LastInsertId()
	a := Article{int(id), title, content, datetime, userID, ""}
	return &a, nil
}
