package main

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func formatStatus(song *Song) string {
	return fmt.Sprintf("Today's No. %s: %s - %s %s %s", song.Hour, song.Band, song.Title, AppConfig.Twitter.Hashtags, song.Link)
}

// UpdateStatusWithSong updates the timeline with a new tweet
func UpdateStatusWithSong(song *Song) error {
	config := oauth1.NewConfig(AppConfig.Twitter.ConsumerKey, AppConfig.Twitter.ConsumerSecret)
	token := oauth1.NewToken(AppConfig.Twitter.AccessToken, AppConfig.Twitter.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	status := formatStatus(song)
	_, _, err := client.Statuses.Update(status, nil)
	if err != nil {
		return err
	}
	return nil
}
