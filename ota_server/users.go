package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/copperhead-security/android_ota_manager/database"
)

// GET /users
func Users(w rest.ResponseWriter, r *rest.Request) {
	users := userApi(database.UserList())
	w.WriteJson(apiListResponse(r, users))
}

// POST /authenticate
func Login(w rest.ResponseWriter, r *rest.Request) {
	// r.ParseForm()
	// session, _ := CookieStore.Get(r, "auth")
	// username := r.FormValue("username")
	// password := r.FormValue("password")
	//
	// user, err := database.FindUserByUsername(username)
	// if err == nil {
	// 	err = scrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// 	if err == nil {
	// 		session.Values["userid"] = user.Id
	// 		session.Save(r, w)
	// 		http.Redirect(w, r, "/admin/releases", http.StatusFound)
	// 	} else {
	// 		log.Println("Login failed", err)
	// 		session.AddFlash("Login failed, bad password!")
	// 		session.Save(r, w)
	// 		http.Redirect(w, r, "/login", http.StatusFound)
	// 	}
	// } else {
	// 	session.AddFlash("Login failed, username not found!")
	// 	session.Save(r, w)
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// }
}

func Logout(w rest.ResponseWriter, r *rest.Request) {
	// session, _ := CookieStore.Get(r, "auth")
	// session.Values["userid"] = nil
	// session.Save(r, w)
	// http.Redirect(w, r, "/login", http.StatusFound)
}

// Generate list of users using ApiResource format
func userApi(items []database.User) []ApiResource {
	resources := make([]ApiResource, 0, len(items))
	for _, item := range items {
		resources = append(resources, userResource(item))
	}
	return resources
}

// Return a user using ApiResource format
func userResource(item database.User) ApiResource {
	resource := Resource(item.Id, "User", "/users")
	resource.Attributes = item
	return resource
}
