package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Song struct
type Song struct {
	Date  string
	Hour  string
	Band  string
	Title string
	Link  string
}

// GetNextSong get the song for the current datetime.
func GetNextSong() (*Song, error) {
	pk, sk, err := generateNextSongKey()
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AppConfig.DynamoDB.DefaultRegion),
	})
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(AppConfig.DynamoDB.SongsTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Date": {
				S: aws.String(pk),
			},
			"Hour": {
				S: aws.String(sk),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	song := Song{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &song)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal the record, %v", err)
	}

	if song.Band == "" {
		return nil, fmt.Errorf("Could not find the item, %v, %v", pk, sk)
	}

	return &song, nil
}

func generateNextSongKey() (string, string, error) {
	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		return "", "", err
	}

	now := time.Now().In(location)
	pk := now.Format(AppConfig.DynamoDB.SongsTableKeyDateFormat)
	sk := now.Format(AppConfig.DynamoDB.SongsTableKeyHourFormat)
	return pk, sk, nil
}
