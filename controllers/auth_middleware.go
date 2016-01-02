package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/copperhead/android_ota_manager/models"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var (
	CookieStore   sessions.Store
	R             *render.Render
)

func InitMiddleware(templates string) {
	// Init cookie store
	ota_key := os.Getenv("OTA_COOKIE_KEY")
	CookieStore = sessions.NewCookieStore([]byte(ota_key))

	// Init renderer
	R = render.New(render.Options{
		Directory:       templates,
		Layout:          "layout",
		Extensions:      []string{".html"},
		IndentJSON:      true,
		PrefixJSON:      []byte(")]}',\n"),
		HTMLContentType: "text/html",
		IsDevelopment:   false,
	})
}

func AuthMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		session, err := CookieStore.Get(r, "auth")
		if err != nil {
			log.Println("Error parsing auth cookie", err)
		}
		user_id := session.Values["userid"]
		if user_id != nil {
			_, err = models.FindUser(user_id.(int64))
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				next(w, r)
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}
