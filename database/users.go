package database

import (
	"log"
	"time"

	"github.com/elithrar/simple-scrypt"
)

type User struct {
	Id       int64 `db:"user_id"`
	Created  int64
	Username string
	Password string
}

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (u *User) CreatedAt() time.Time {
	t := time.Unix(0, u.Created)
	return t
}

func (u *User) HumanCreatedAt() string {
	createdAt := u.CreatedAt()
	return createdAt.Format(time.RFC3339)
}

func UserList() []User {
	var users []User
	_, err := dbmap.Select(&users, "select * from users order by user_id DESC")
	if err != nil {
		log.Println("Select all users failed", err)
	}
	return users
}

func FindUser(id int64) (User, error) {
	var user User
	err := dbmap.SelectOne(&user, "select * from users where user_id=? LIMIT 1", id)
	if err != nil {
		log.Println("Find user by user_id failed", id, err)
	}
	return user, err
}

func FindUserByUsername(username string) (User, error) {
	var user User
	err := dbmap.SelectOne(&user, "select * from users where username=? LIMIT 1", username)
	if err != nil {
		log.Println("Find user by username failed", username, err)
	}
	return user, err
}

func CreateUser(username string, password string) error {
	hash, err := scrypt.GenerateFromPassword([]byte(password), scrypt.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	user := User{
		Created:  time.Now().UnixNano(),
		Username: username,
		Password: string(hash),
	}

	// Insert into db
	err = dbmap.Insert(&user)
	if err != nil {
		log.Println("Create user failed", username, err)
	}
	return err
}
