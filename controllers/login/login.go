package login

import (
	"net/http"
	"text/template"
	//"github.com/gorilla/sessions"
)

// var (
// 	key   = []byte(config.SESSION_KEY)
// 	store = sessions.NewCookieStore(key)
// )

func Login(w http.ResponseWriter, r *http.Request) {
	// session, err := store.Get(r, "user_name")
	// if err != nil {
	// 	//redirect to main page
	// }

	viewPage := "views/login.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	// if r.Method == "POST" {

	//dynamo DB

	// 	return
	// }

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets":            assetsUrl,
		"username_err_bool": false,
		"password_err_bool": false,
		"username_err":      "",
		"password_err":      "",
		"login_err":         "",
		"login_err_bool":    false,
		"username_filled":   false,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "login", data)
}
