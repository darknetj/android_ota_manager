package models

import (
    "log"
)

type User struct {
    Id          int64 `db:"user_id"`
    Created     int64
    Username    string
    Password    string
}

type LoginForm struct {
    User     string `form:"user" binding:"required"`
    Password string `form:"password" binding:"required"`
}

func UserList() []User {
    var users []User
    _, err := dbmap.Select(&users, "select * from users order by user_id DESC")
    if err != nil {
      log.Println("Select all users failed", err)
    }
    return users
}

func FindUser(user User, id int64) error {
    err := dbmap.SelectOne(&user, "select * from users where user_id=?", id)
    if err != nil {
      log.Println("Find user by user_id failed", id, err)
    }
    return err
}

func FindUserByUsername(user User, username string) error {
    err := dbmap.SelectOne(&user, "select * from users where username=?", username)
    if err != nil {
      log.Println("Find user by username failed", username, err)
    }
    return err
}
