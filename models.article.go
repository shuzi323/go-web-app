package main

import (
	"errors"
	"log"
	"strings"
)

type article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func getAllArticles() (articles []article, err error) {
	articles = make([]article, 0)
	rows, err := db.Query("SELECT article_id, title, content FROM article ORDER BY article_id DESC")
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var perArticle article
		rows.Scan(&perArticle.ID, &perArticle.Title, &perArticle.Content)
		articles = append(articles, perArticle)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func getArticleByID(id int) (*article, error) {
	var perArticle article
	err := db.QueryRow("SELECT title, content FROM article WHERE article_id=?", id).Scan(&perArticle.Title, &perArticle.Content)
	return &perArticle, err
}

func createNewArticle(title, content string) (*article, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("Title can'n be empty")
	}
	stmt, err := db.Exec("INSERT INTO article(title, content) VALUES(?, ?)", title, content)
	if err != nil {
		return nil, err
	}
	id, err := stmt.LastInsertId()
	a := article{int(id), title, content}
	return &a, nil
}
