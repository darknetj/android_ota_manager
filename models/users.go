package models

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
