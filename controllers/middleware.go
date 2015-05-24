package controllers

import (
    "os"
    "log"
    "net/http"
    "github.com/unrolled/render"
    "github.com/gorilla/sessions"
    "github.com/copperhead-security/android_ota_server/models"
)

var (
  CookieStore sessions.Store
  R *render.Render
)

func InitMiddleware() {
  // Init cookie store
  ota_key := os.Getenv("OTA_COOKIE_KEY")
  CookieStore = sessions.NewCookieStore([]byte(ota_key))

  // Init renderer
  R = render.New(render.Options{
    Directory: "views",
    Layout: "layout",
    Extensions: []string{".html"},
    IndentJSON: true,
    PrefixJSON: []byte(")]}',\n"),
    HTMLContentType: "text/html",
    IsDevelopment: false,
  })
}

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
