package main

import (
	"errors"
	"log"
	"strings"
)

type User struct {
	ID       int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"-"`
}

//验证用户 并返回ID
func isUserValid(username, password string) (bool, int) {
	var userQuery User
	err := db.QueryRow("SELECT user_id, username, password FROM user WHERE username = ?", username).Scan(&userQuery.ID, &userQuery.Username, &userQuery.Password)
	if err != nil {
		log.Println(err)
		return false, -1
	}
	if userQuery.Username == username && userQuery.Password == password {
		return true, userQuery.ID
	}
	return false, -1
}

func isUserNameAvailable(username string) bool {
	var userQuery User
	err := db.QueryRow("SELECT username, password FROM user WHERE username = ?", username).Scan(&userQuery.Username, &userQuery.Password)
	if err != nil {
		log.Println(err) //注册用户ID不存在，则sql: no rows in result set
		return true
	}
	return false
}

func registerNewUser(username, password string) (*User, error) {
	log.Println(1)
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUserNameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}
	log.Println(2)
	u := User{Username: username, Password: password}
	//stmt 可用于获取用户ID，这里没必要，所以省略
	stmt, err := db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		log.Println(err)
	}
	id, err := stmt.LastInsertId()
	u.ID = int(id)
	return &u, nil
}

//保存token到数据库
func storeToken(token string, id int) {
	_, err := db.Exec("INSERT INTO token(token_data, user_id) VALUES(?, ?) ON DUPLICATE KEY UPDATE token_data=?", token, id, token)
	if err != nil {
		log.Println(err)
	}
}
