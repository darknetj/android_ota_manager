package controllers

import (
    "log"
    "net/http"
    // "strconv"
    "github.com/elithrar/simple-scrypt"
    "github.com/copperhead-security/android_ota_server/models"
    "github.com/copperhead-security/android_ota_server/lib"
)

// GET /users
func Users(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {"users": models.UserList()}
    lib.T("users.html").Execute(w, data)
}

// POST /authenticate
func Authenticate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    log.Println("auuth")
    username := r.FormValue("username")
    password := r.FormValue("password")

    var user models.User
    user = models.FindUserByUsername(username)

    // Uses the parameters from the existing derived key. Return an error if they don't match.
    err := scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        log.Println(err)
        http.Redirect(w, r, "/login", http.StatusFound)
    } else {
        // session := sessions.Get(c)
        // log.Println(session)
        // log.Println(user.Id)
        // session.Set("userid", strconv.FormatInt(user.Id,10))
        // session.Save()
        http.Redirect(w, r, "/releases", http.StatusFound)
    }
}

func Login(w http.ResponseWriter, r *http.Request) {
}

func Logout(w http.ResponseWriter, r *http.Request) {
}
