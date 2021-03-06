package controllers

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/eoria17/AWS-Golang-Music-Sub/models"
)

func (ae AppEngine) Register(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user_cookie")
	if err != nil {
		fmt.Println(err)
	}

	//if logged in
	if auth, ok := session.Values["logged_in"].(bool); ok || auth {
		http.Redirect(w, r, "http://"+r.Host+"/home", http.StatusPermanentRedirect)
		return
	}

	//variable declarations
	viewPage := "views/register.html"
	assetsUrl := "http://" + r.Host + "/assets/"

	username_err, email_err, password_err := "", "", ""
	username_err_bool, password_err_bool, email_err_bool := false, false, false
	username_filled, email_filled := false, false

	username, email := "", ""

	if r.Method == "POST" {
		r.ParseForm()

		//checks if empty
		if r.FormValue("email") == "" {
			email_err_bool = true
			email_err = "Please enter email."
		} else {
			email_filled = true
			email = r.FormValue("email")
		}

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

		//search DB for register data
		if !username_err_bool && !password_err_bool && !email_err_bool {
			user := ae.GetCurrentUser(r.FormValue("email"))

			if user.Email != "" {
				email_err_bool = true
				email_err = "Email already exist, please enter a different email."
			} else {

				newUser := models.Login{
					Email:     r.FormValue("email"),
					User_name: r.FormValue("username"),
					Password:  r.FormValue("password"),
				}

				av, err := dynamodbattribute.MarshalMap(newUser)
				if err != nil {
					log.Fatalf("Got error marshalling new movie item: %s", err)
				}

				input := &dynamodb.PutItemInput{
					Item:      av,
					TableName: aws.String("login"),
				}

				_, err = ae.DynamoDBClient.PutItem(input)
				if err != nil {
					log.Fatalf("Got error calling PutItem: %s", err)
				}

				//redirect to login
				http.Redirect(w, r, "http://"+r.Host+"/", http.StatusPermanentRedirect)
				return
			}
		}

	}

	t, _ := template.ParseFiles(viewPage)

	data := map[string]interface{}{
		"assets":            assetsUrl,
		"username_err_bool": username_err_bool,
		"password_err_bool": password_err_bool,
		"email_err_bool":    email_err_bool,
		"email_err":         email_err,
		"username_err":      username_err,
		"password_err":      password_err,
		"username_filled":   username_filled,
		"username":          username,
		"email_filled":      email_filled,
		"email":             email,
	}

	w.WriteHeader(http.StatusOK)
	t.ExecuteTemplate(w, "register", data)
}
