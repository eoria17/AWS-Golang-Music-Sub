package controllers

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/eoria17/AWS-Golang-Music-Sub/models"
	//"github.com/gorilla/sessions"
)

// var (
// 	key   = []byte(config.SESSION_KEY)
// 	store = sessions.NewCookieStore(key)
// )

func (ae AppEngine) Login(w http.ResponseWriter, r *http.Request) {

	// session, err := store.Get(r, "user_name")
	// if err != nil {
	// 	//redirect to main page
	// }

	viewPage := "views/login.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	login_err := ""
	password_err := ""
	username_err := ""
	username_err_bool := false
	password_err_bool := false
	login_err_bool := false
	username_filled := false
	username := ""

	if r.Method == "POST" {

		r.ParseForm()

		//check if username or password is null
		if r.FormValue("username") == "" {
			username_err_bool = true
			username_err = "Please enter username."
		} else {
			username_filled = true
			username = r.FormValue("username")
		}

		if r.FormValue("password") == "" {
			password_err_bool = true
			password_err = "Please enter password."
		}

		//search DB for login data
		if !username_err_bool && !password_err_bool {
			svc := dynamodb.New(ae.Session)

			result, err := svc.GetItem(&dynamodb.GetItemInput{
				TableName: aws.String("login"),
				Key: map[string]*dynamodb.AttributeValue{
					"user_name": {
						S: aws.String(r.FormValue("username")),
					},
				},
			})

			if err != nil {
				fmt.Println("here", err)
			}

			if result.Item == nil {
				login_err = "invalid username or password"
			}

			user := models.Login{}
			err = dynamodbattribute.UnmarshalMap(result.Item, &user)

			if err != nil {
				panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
			}

			if user.Password != r.FormValue("password") {
				login_err = "invalid username or password"
				login_err_bool = true
			} else if user.Password == r.FormValue("password") {
				fmt.Println("meep")
				//redirect to home page
			}

		}

	}

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets":            assetsUrl,
		"username_err_bool": username_err_bool,
		"password_err_bool": password_err_bool,
		"username_err":      username_err,
		"password_err":      password_err,
		"login_err":         login_err,
		"login_err_bool":    login_err_bool,
		"username_filled":   username_filled,
		"username":          username,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "login", data)
}
