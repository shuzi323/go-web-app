package main

import (
	"errors"
	"log"
	"strings"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

//验证用户
func isUserValid(username, password string) bool {
	var userQuery user
	err := db.QueryRow("SELECT username, password FROM user WHERE username = ?", username).Scan(&userQuery.Username, &userQuery.Password)
	if err != nil {
		log.Println(err)
		return false
	}
	if userQuery.Username == username && userQuery.Password == password {
		return true
	}
	return false
}

func isUserNameAvailable(username string) bool {
	var userQuery user
	err := db.QueryRow("SELECT username, password FROM user WHERE username = ?", username).Scan(&userQuery.Username, &userQuery.Password)
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func registerNewUser(username, password string) (*user, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("The password can't be empty")
	} else if !isUserNameAvailable(username) {
		return nil, errors.New("The username isn't available")
	}

	u := user{Username: username, Password: password}
	//stmt 可用于获取用户ID，这里没必要，所以省略
	_, err := db.Exec("INSERT INTO user(username, password) VALUES(?, ?)", username, password)
	if err != nil {
		log.Println(err)
	}
	return &u, nil
}
