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
}

func getAllArticles() (articles []Article, err error) {
	articles = make([]Article, 0)
	rows, err := db.Query("SELECT article_id, title, content, datetime FROM article ORDER BY article_id DESC")
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var perArticle Article
		rows.Scan(&perArticle.ID, &perArticle.Title, &perArticle.Content, &perArticle.DateTime)
		articles = append(articles, perArticle)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func getArticleByID(id int) (*Article, error) {
	var perArticle Article
	err := db.QueryRow("SELECT title, content, datetime FROM article WHERE article_id=?", id).Scan(&perArticle.Title, &perArticle.Content, &perArticle.DateTime)
	return &perArticle, err
}

func createNewArticle(title, content string) (*Article, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("Title can'n be empty")
	}
	datetime := time.Now().Local().Format("2006-01-02 15:04:05")
	stmt, err := db.Exec("INSERT INTO article(title, content, datetime) VALUES(?, ?, ?)", title, content, datetime)
	if err != nil {
		return nil, err
	}
	id, err := stmt.LastInsertId()
	a := Article{int(id), title, content, datetime}
	return &a, nil
}
