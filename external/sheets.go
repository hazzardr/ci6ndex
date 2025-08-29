package external

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	credentialFileLocation string = "token.json"
	readOnlySheetsScope    string = "https://www.googleapis.com/auth/spreadsheets.readonly"
)

type GoogleSheets struct {
	client *sheets.Service
}

func NewGoogleSheets(oauthConfigFileLocation string) (*GoogleSheets, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := buildAuthenticatedClient(ctx, oauthConfigFileLocation)
	if err != nil {
		return nil, err
	}

	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to build google sheets client: %w", err)
	}

	return &GoogleSheets{
		client: service,
	}, nil
}

func buildAuthenticatedClient(ctx context.Context, oauthConfigFileLocation string) (*http.Client, error) {
	cf, err := os.ReadFile(oauthConfigFileLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to read oauth2 info from disk: %w", err)
	}
	config, err := google.ConfigFromJSON(cf, readOnlySheetsScope)
	if err != nil {
		return nil, err
	}
	token, err := tokenFromFile(credentialFileLocation)
	if err != nil {
		slog.Warn("failed to find cached access token, trying web...", "err", err)
		token, err = tokenFromWeb(ctx, config)
		if err != nil {
			return nil, err
		}
	}
	err = tokenToFile(credentialFileLocation, token)
	if err != nil {
		slog.Warn("failed to cache access token to local disk, continuing...", "err", err)
	}

	client := config.Client(ctx, token)
	return client, nil
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("please navigate to the following URL and then type in the access code: \n%v\n", authURL)
	var authCode string
	_, err := fmt.Scan(&authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to read authorization code from console: %w", err)
	}

	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve token from google: %w", err)
	}
	return token, nil
}

func tokenFromFile(filepath string) (*oauth2.Token, error) {
	secret, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open client secret file: %w", err)
	}
	defer secret.Close()
	token := &oauth2.Token{}
	err = json.UnmarshalRead(secret, token)
	return token, err
}

func tokenToFile(filepath string, token *oauth2.Token) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to cache token to file: %w", err)
	}
	defer f.Close()
	err = json.MarshalWrite(f, token)
	return err
}
