package api

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log/slog"
	"net/http"
	"os"
)

func getRankings(c *AppConfig) error {
	ctx := context.Background()
	oauthConfigFile, err := os.ReadFile(c.GoogleCloudCredentialsLocation)
	if err != nil {
		return err
	}
	oauthConfig, err := google.ConfigFromJSON(oauthConfigFile, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return err
	}

	client := getClient(oauthConfig)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	// Just try to get connection for now
	spreadsheetId := "19fO_8daYFEdCqla1YwSSaCfk2zEq7SleKRvIUbt36Jg"

}

// Retrieve a token, save the token, return generated client.
// From quickstart: https://developers.google.com/sheets/api/quickstart/go
func getClient(config *oauth2.Config) *http.Client {
	tFile := "token.json"
	t, err := getTokenFromFile(tFile)
	if err != nil {
		slog.Warn("could not retrieve gcloud auth from file, trying web")
		t = getTokenFromWeb(config)
		saveToken(tFile, t)
	}
	return config.Client(context.Background(), t)
}

func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	return t, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	slog.Info("go to the link and type authorization code from there", "authUrl", authUrl)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		slog.Error("unable to read authorization code from prompt", "error", err)
		os.Exit(1)
	}

	t, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		// TODO: put this behind rest and just give bad status instead of crash
		slog.Error("unable to retrieve google sheets access token", "error", err)
		os.Exit(1)
	}
	return t
}

// Saves the token to a file path
func saveToken(path string, token *oauth2.Token) {
	slog.Info("saving gcloud credential file", "path", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		// TODO: put this behind rest and just give bad status instead of crash
		slog.Error("unable to save google sheets access token to file", "error", err)
		os.Exit(1)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
