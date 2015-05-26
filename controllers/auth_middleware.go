package controllers

import (
    "os"
    "log"
    "net/http"
    "github.com/unrolled/render"
    "github.com/gorilla/sessions"
    "github.com/codegangsta/negroni"
    "github.com/copperhead-security/android_ota_server/models"
)

var (
  CookieStore sessions.Store
  R *render.Render
  TemplatesPath string
)

func InitMiddleware(templates string) {
  // Init cookie store
  ota_key := os.Getenv("OTA_COOKIE_KEY")
  CookieStore = sessions.NewCookieStore([]byte(ota_key))

  // Init renderer
  R = render.New(render.Options{
    Directory: templates,
    Layout: "layout",
    Extensions: []string{".html"},
    IndentJSON: true,
    PrefixJSON: []byte(")]}',\n"),
    HTMLContentType: "text/html",
    IsDevelopment: false,
  })
}

func AuthMiddleware() negroni.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    session, err := CookieStore.Get(r, "auth")
    if err != nil {
      log.Println("Error parsing auth cookie", err)
    }
    var user models.User
    user_id := session.Values["userid"]
    if user_id != nil {
      user, err = models.FindUser(user_id.(int64))
      if err != nil {
          http.Redirect(w, r, "/login", http.StatusFound)
      } else {
          log.Println(user)
          next(w, r)
      }
    } else {
        http.Redirect(w, r, "/login", http.StatusFound)
    }
  }
}
