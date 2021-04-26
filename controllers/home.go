package controllers

import (
	"fmt"
	"net/http"
	"text/template"
)

func (ae AppEngine) Home(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if logged in
	if auth, ok := session.Values["logged_in"].(bool); !ok || !auth {
		http.Redirect(w, r, "http://"+r.Host+"/", http.StatusPermanentRedirect)
		return
	}

	viewPage := "views/home.html"
	assetsUrl := "http://" + r.Host + "/assets/"
	user := ae.GetCurrentUser(session.Values["email"].(string))

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets":   assetsUrl,
		"username": user.User_name,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "home", data)
}
