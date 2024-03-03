package internal

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
	"slices"
	"strconv"
)

type Ranking struct {
	Player            string
	CombinedLeaderCiv string
	Tier              float64
}

type SheetsServiceProvider interface {
	GetClient(ctx context.Context, credsLocation string) (*sheets.Service, error)
}

type SheetsService struct {
}

func (s *SheetsService) GetClient(ctx context.Context, credsLocation string) (*sheets.Service, error) {
	oauthConfigFile, err := os.ReadFile(credsLocation)
	if err != nil {
		return nil, err
	}
	oauthConfig, err := google.ConfigFromJSON(oauthConfigFile, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}

	client, err := getClient(oauthConfig)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

// Get token from file
// if it's not in the file / file DNE then get from the web
// create http client from that token
// create sheets service from that client
// get the spreadsheet
// if there is auth error, refresh token from web + save to file
// retry get spreadsheet
// parse spreadsheet

// OAuth2TokenProvider gets an OAuth2 token from file or web
type OAuth2TokenProvider interface {
	GetTokenFromFile(file string) (*oauth2.Token, error)
	GetTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error)
}

type OAuth2TokenService struct {
	GetTokenFromFile func(file string) (*oauth2.Token, error)
	GetTokenFromWeb  func(config *oauth2.Config) (*oauth2.Token, error)
}

func GetRankingsFromSheets(config *AppConfig, ctx context.Context) ([]Ranking, error) {
	sheetService := &SheetsService{}
	srv, err := sheetService.GetClient(ctx, config.GoogleCloudCredentialsLocation)
	if err != nil {
		return nil, err
	}

	civLeaderCells := "A2:K80"
	ss, err := srv.Spreadsheets.Values.Get(config.CivRankingSheetId, civLeaderCells).MajorDimension("ROWS").Do()
	if err != nil {
		var tokenErr *oauth2.RetrieveError
		if errors.As(err, &tokenErr) {
			if tokenErr.ErrorCode == "invalid_grant" {
				slog.Warn("could not retrieve token from file, attempting to refresh...",
					"error", tokenErr)
				t, err := getTokenFromWeb(oauth2Config)
				if err != nil {
					slog.Error("could not refresh token", "error", err)
					return nil, err
				}
			}
			return nil, tokenErr

		}
		return nil, err
	}
	if ss.ServerResponse.HTTPStatusCode != 200 {
		slog.Error("could not get spreadsheet", "status", ss.ServerResponse.HTTPStatusCode)
		return nil, errors.New("could not get spreadsheet")
	}

	var rankings []Ranking
	ignoreNames := []string{"Civ", "Average", "Deviation"}

	for rowNum, row := range ss.Values {
		r := Ranking{}
		for colNum, cell := range row {
			// Labels
			if rowNum == 0 || colNum == 0 {
				continue
			}

			ranker, ok := ss.Values[0][colNum].(string)
			if !ok {
				slog.Error("sheets header data is malformed, stopping processing", "colNum", colNum)
				return nil, fmt.Errorf("sheets header data is malformed at column=%disc, stopping processing", colNum)
			}

			// Ignore computed fields
			if slices.Contains(ignoreNames, ranker) {
				continue
			}

			val, ok := cell.(string)
			if !ok {
				slog.Warn("found non-string in sheets tier cell, skipping", "colNum", colNum, "rowNum", rowNum, "cell", cell)
				continue
			}

			if val == "" {
				continue
			}

			tier, err := strconv.ParseFloat(val, 64)
			if err != nil {
				slog.Warn("could not parse tier from string, skipping", "string", val, "error", err)
				continue
			}
			r.Tier = tier
			r.Player = ranker
			r.CombinedLeaderCiv = ss.Values[rowNum][0].(string)
			rankings = append(rankings, r)
		}
	}

	return rankings, nil

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
