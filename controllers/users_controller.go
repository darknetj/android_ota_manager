package controllers

import (
    "log"
    "os"
    "net/http"
    // "strconv"
    "github.com/gorilla/sessions"
    "github.com/elithrar/simple-scrypt"
    "github.com/copperhead-security/android_ota_server/models"
    "github.com/copperhead-security/android_ota_server/lib"
)

var (
  store sessions.Store
)

// GET /users
func Users(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {"users": models.UserList()}
    lib.T("users.html").Execute(w, data)
}

// POST /authenticate
func Authenticate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    session, _ := store.Get(r, "auth")
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
    session, _ := store.Get(r, "auth")
    data := map[string]interface{} {
      "flashes": session.Flashes(),
    }
    session.Save(r, w)
    lib.T("user_login.html").Execute(w, data)
}

func Logout(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/login", http.StatusFound)
}

func InitAuth() {
  // Init cookie store
  ota_key := os.Getenv("OTA_COOKIE_KEY")
  store = sessions.NewCookieStore([]byte(ota_key))
}

// Middleware to check auth
func Auth(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // session := sessions.Get(c)
        var user models.User
        user_id := int64(1)
        err := models.FindUser(user, user_id)
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusFound)
        } else {
            log.Println(user)
            next.ServeHTTP(w, r)
        }
        // defer func() {
        //     if r := recover(); r != nil {
        //         c.String(http.StatusInternalServerError, "Unauthorized!")
        //         c.Fail(401, errors.New("Unauthorized"))
        //     }
        // }()
        // user_id := session.Get("userid")
        // switch user_id.(type) {
        //     case string:
        //         user = models.FindUser(user_id)
        //         log.Println(user)
        //         if user != nil {
        //             /* panic("Not authenticated") */
        //         } else {
        //             next.ServeHTTP(w, r)
        //         }
        //     default:
        //         /* panic("Not authenticated") */
        // }
    })
}
