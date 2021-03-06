package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eoria17/AWS-Golang-Music-Sub/controllers"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/eoria17/AWS-Golang-Music-Sub/config"
)

func main() {
	//create AWS session
	creds := credentials.NewStaticCredentials(config.ACCESS_KEY_ID, config.SECRET_ACCESS_KEY, "")
	creds.Get()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: creds,
	})

	if err != nil {
		fmt.Println(err)
	}

	//create router handler
	router := mux.NewRouter()

	//dynamoDB AWS client
	svc := dynamodb.New(sess)

	//S3 AWS client
	//s3c := s3.New(sess)
	s3uploader := s3manager.NewUploader(sess)

	//dependency injection
	appEngine := controllers.AppEngine{
		DynamoDBClient: svc,
		S3Client:       s3uploader,
	}

	//routing
	appEngine.Route(router)
	appEngine.DataSeed()

	//serve public as static file
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/"))))
	http.Handle("/assets/", router)

	//run server
	fmt.Println("Currently Listening to port 8080..")
	log.Println(http.ListenAndServe(":8080", router))

}
