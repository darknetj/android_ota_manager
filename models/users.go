package models

import (
    "log"
    "github.com/copperhead-security/android_ota_server/lib"
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
    lib.CheckErr(err, "Select all users failed")
    return users
}

func FindUser(id int64) User {
    var user User
    err := dbmap.SelectOne(&user, "select * from users where user_id=?", id)
    if err != nil {
      log.Println(err)
    }
    return user
}

func FindUserByUsername(username string) User {
    var user User
    err := dbmap.SelectOne(&user, "select * from users where username=?", username)
    if err != nil {
      log.Println(err)
    }
    return user
}
