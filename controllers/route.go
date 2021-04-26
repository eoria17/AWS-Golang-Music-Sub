package controllers

import (
	"fmt"

	"github.com/eoria17/AWS-Golang-Music-Sub/models"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type AppEngine struct {
	Session        *session.Session
	DynamoDBClient *dynamodb.DynamoDB
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
	r.HandleFunc("/home", ae.Home)
	r.HandleFunc("/register", ae.Register)
}

func (ae AppEngine) GetCurrentUser(username string) models.Login {
	svc := ae.DynamoDBClient

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("login"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(username),
			},
		},
	})

	if err != nil {
		fmt.Println("here", err)
	}

	user := models.Login{}

	if result != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &user)

		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
	}

	return user
}
