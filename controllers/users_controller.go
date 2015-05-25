package controllers

import (
    "log"
    "net/http"
    "strconv"
    "github.com/elithrar/simple-scrypt"
    "github.com/copperhead-security/android_ota_server/models"
)

// GET /users
func Users(w http.ResponseWriter, r *http.Request) {
    users := models.UserList()
    var userStrings []map[string]string
    for index, user := range users {
      log.Println(index, user)
      userStrings = append(userStrings, map[string]string{
        "Id": strconv.FormatInt(user.Id, 10),
        "Username": user.Username,
        "Created": user.HumanCreatedAt(),
      })
    }
    data := map[string][]models.User { "users": models.UserList()}
    log.Println(data)
    R.HTML(w, http.StatusOK, "users", data)
}

// POST /authenticate
func Authenticate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    session, _ := CookieStore.Get(r, "auth")
    username := r.FormValue("username")
    password := r.FormValue("password")

    user, err := models.FindUserByUsername(username)
    if err == nil {
        log.Println(user)
        log.Println(user.Password)
        err = scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
        if err == nil {
            session.Values["userid"] = user.Id
            session.Save(r, w)
            http.Redirect(w, r, "/admin/releases", http.StatusFound)
        } else {
            log.Println("Login failed", err)
            session.AddFlash("Login failed, bad password!")
            session.Save(r, w)
            http.Redirect(w, r, "/login", http.StatusFound)
        }
    } else {
        session.AddFlash("Login failed, username not found!")
        session.Save(r, w)
        http.Redirect(w, r, "/login", http.StatusFound)
    }
}

func Login(w http.ResponseWriter, r *http.Request) {
    session, _ := CookieStore.Get(r, "auth")
    data := map[string]interface{} {
      "flashes": session.Flashes(),
      "noAuth": true,
    }
    session.Save(r, w)
    R.HTML(w, http.StatusOK, "user_login", data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    session, _ := CookieStore.Get(r, "auth")
    session.Values["userid"] = nil
    session.Save(r, w)
    http.Redirect(w, r, "/login", http.StatusFound)
}
