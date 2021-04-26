package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/eoria17/AWS-Golang-Music-Sub/config"
	"github.com/eoria17/AWS-Golang-Music-Sub/models"
	"github.com/gorilla/mux"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AppEngine struct {
	DynamoDBClient *dynamodb.DynamoDB
	S3Client       *s3manager.Uploader
}

func (ae AppEngine) Route(r *mux.Router) {
	r.HandleFunc("/", ae.Login)
	r.HandleFunc("/main", ae.Main)
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

func (ae AppEngine) DataSeed() {
	fmt.Println("Initiating data seed..")
	CreateSongTable(ae)
	InsertSongs(ae)
	UploadImages(ae)
}

func UploadImages(ae AppEngine) {

	songs := GetSongsFromJSON()

	for _, song := range songs {
		_, err := ae.S3Client.Upload(&s3manager.UploadInput{
			Bucket: aws.String(config.BUCKET_NAME),
			Key:    aws.String(path.Base(song.ImgURL)),
			Body:   bytes.NewReader(DownloadImage(song.ImgURL)),
			ACL:    aws.String("public-read"),
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}

func DownloadImage(URL string) []byte {
	response, err := http.Get(URL)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Received non 200 response code")
		return nil
	}

	image, _, err := image.Decode(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	buffer := new(bytes.Buffer)
	jpeg.Encode(buffer, image, nil)

	return buffer.Bytes()
}

func CreateSongTable(ae AppEngine) {
	fmt.Println("Creating table : music")
	tableName := "music"

	//create table
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("title"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
		TableName: aws.String(tableName),
	}

	_, err := ae.DynamoDBClient.CreateTable(input)
	if err != nil {
		fmt.Println("table already exist, table will not be created")
	}

}

func GetSongsFromJSON() []models.Song {
	raw, err := ioutil.ReadFile("./a2.json")
	if err != nil {
		log.Fatalf("Got error reading file: %s", err)
	}

	a2 := models.A2{}
	json.Unmarshal(raw, &a2)
	return a2.Songs
}

func InsertSongs(ae AppEngine) {
	fmt.Println("Importing songs from a2.json")
	songs := GetSongsFromJSON()

	for _, song := range songs {

		av, err := dynamodbattribute.MarshalMap(song)
		if err != nil {
			log.Fatalf("Got error marshalling new movie item: %s", err)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String("music"),
		}

		_, err = ae.DynamoDBClient.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}
	}
}
