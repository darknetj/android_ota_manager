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
    log.Println("auuth")
    username := r.FormValue("username")
    password := r.FormValue("password")

    var user models.User
    err := models.FindUserByUsername(user, username)
    if err != nil {
        // Uses the parameters from the existing derived key. Return an error if they don't match.
        err = scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
        if err != nil {
            log.Println("Login failed", err)
            session.AddFlash("Login failed, bad password!")
            http.Redirect(w, r, "/login", http.StatusFound)
        } else {
            // session := sessions.Get(c)
            // log.Println(session)
            // log.Println(user.Id)
            // session.Set("userid", strconv.FormatInt(user.Id,10))
            // session.Save()
            http.Redirect(w, r, "/releases", http.StatusFound)
        }
    } else {
        session.AddFlash("Login failed, username not found!")
        http.Redirect(w, r, "/login", http.StatusFound)
    }
}

func Login(w http.ResponseWriter, r *http.Request) {
    session, _ := CookieStore.Get(r, "auth")
    data := map[string]interface{} {
      "flashes": session.Flashes(),
    }
    session.Save(r, w)
    R.HTML(w, http.StatusOK, "user_login", data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/login", http.StatusFound)
}
