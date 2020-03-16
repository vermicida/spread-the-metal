package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vermicida/spread-the-metal/pkg/config"
)

// CloudWatchEvent struct represent an event from CloudWatch
type CloudWatchEvent struct{}

// AppConfig is the application configuration
var AppConfig = config.New()

// HandleRequest is the lambda function handler
func HandleRequest(ctx context.Context, event CloudWatchEvent) (string, error) {
	song, err := GetNextSong()
	if err != nil {
		return "Something went wrong retrieving the next song", err
	}

	err = UpdateStatusWithSong(song)
	if err != nil {
		return "Something went wrong updating the Twitter status", err
	}

	return fmt.Sprintf("%s's %s has been tweeted :-)", song.Band, song.Title), nil
}

func main() {
	lambda.Start(HandleRequest)
}
