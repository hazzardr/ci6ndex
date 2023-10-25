package api

import (
	"context"
	"encoding/json"
	"errors"
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

	client, err := getClient(oauthConfig)
	if err != nil {
		return err
	}
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	civLeaderCells := "A2:K80"
	ss, err := srv.Spreadsheets.Values.BatchGet(c.CivRankingSheetId).Ranges(civLeaderCells).Do()
	//ss, err := srv.Spreadsheets.Values.Get(c.CivRankingSheetId, readRange).Do()
	if err != nil {
		return err
	}
	if ss.ServerResponse.HTTPStatusCode != 200 {
		slog.Error("could not get spreadsheet", "status", ss.ServerResponse.HTTPStatusCode)
		return errors.New("could not get spreadsheet")
	}

	return nil

}

// Retrieve a token, save the token, return generated client.
// From quickstart: https://developers.google.com/sheets/api/quickstart/go
func getClient(config *oauth2.Config) (*http.Client, error) {
	tFile := "token.json"
	t, err := getTokenFromFile(tFile)
	if err != nil {
		slog.Warn("could not retrieve gcloud auth from local file, trying web")
		t, err = getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		err = saveToken(tFile, t)
		if err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), t), nil
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

func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	slog.Info("go to the link and type authorization code from there", "authUrl", authUrl)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		slog.Error("unable to read authorization code from prompt", "error", err)
		return nil, err
	}

	t, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		slog.Error("unable to retrieve google sheets access token", "error", err)
		return nil, err
	}
	return t, nil
}

// Saves the token to a file path
func saveToken(path string, token *oauth2.Token) error {
	slog.Info("saving gcloud credential file", "path", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		// TODO: put this behind rest and just give bad status instead of crash
		slog.Error("unable to save google sheets access token to file", "error", err)
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return err
	}
	return nil
}
