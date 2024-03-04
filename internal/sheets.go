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
	"os"
	"strconv"
	"time"
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
	token *oauth2.Token
}

func (s *SheetsService) GetClient(ctx context.Context, oauthConfigLoc string) (*sheets.Service,
	error) {

	oauthConfig, err := getOAuthConfigFromFile(oauthConfigLoc)
	if err != nil {
		return nil, err
	}

	client := oauthConfig.Client(ctx, s.token)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

// OAuth2TokenProvider gets an OAuth2 token from file or web
type OAuth2TokenProvider interface {
	GetToken() (*oauth2.Token, error)
}

type OAuth2TokenService struct {
	oauthCredsLocation string
	tokenFileLocation  string
}

// GetToken retrieves a valid OAuth2 token for interacting with Google Sheets.
// It will try to:
//  1. Get the token from a local file
//  2. If the file does not exist, is malformed, or is expired,
//     it will call out to the web for a new token.
func (o *OAuth2TokenService) GetToken() (*oauth2.Token, error) {
	t, err := getTokenFromFile(o.tokenFileLocation)
	if err != nil {
		slog.Warn("could not retrieve oauth2 token from local file, trying web")
		t, webErr := getTokenFromWeb(o.oauthCredsLocation)
		if webErr != nil {
			slog.Error("could not retrieve oauth2 token from file or web",
				"fileErr", err, "webErr", webErr)
			return nil, errors.Join(err, webErr)
		}

		err = saveToken(o.tokenFileLocation, t)
		if err != nil {
			slog.Warn("retrieved token from web, but could not save token to file", "error", err)
			return t, nil
		}
	}
	if t.Expiry.Before(time.Now()) {
		slog.Info("retrieved token has expired. refreshing...")
		t, err = getTokenFromWeb(o.oauthCredsLocation)
		if err != nil {
			slog.Error("could not refresh token from web", "error", err)
			return nil, err
		}
	}
	return t, nil
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

func GetRankingsFromSheets(config *AppConfig, ctx context.Context) ([]Ranking, error) {
	ts := &OAuth2TokenService{
		oauthCredsLocation: config.GoogleCloudCredentialsLocation,
		tokenFileLocation:  "token.json",
	}

	t, err := ts.GetToken()
	if err != nil {
		return nil, err
	}

	sheetService := &SheetsService{
		token: t,
	}

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
				t, err := ts.GetToken()
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

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
// It uses oauth2 config stored at oauthConfigLocation as a readable file.
func getTokenFromWeb(oauthConfigLocation string) (*oauth2.Token, error) {
	oauthConfig, err := getOAuthConfigFromFile(oauthConfigLocation)
	if err != nil {
		return nil, errors.Join(errors.New("could not read oauth2 file config"), err)
	}

	authUrl := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	slog.Info("go to the following link and type authorization code from there in this terminal. "+
		"Thanks!", "authUrl", authUrl)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		slog.Error("unable to read authorization code from prompt", "error", err)
		return nil, err
	}

	t, err := oauthConfig.Exchange(context.Background(), authCode)
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

func getOAuthConfigFromFile(fileLocation string) (*oauth2.Config, error) {
	f, err := os.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}

	oauthConfig, err := google.ConfigFromJSON(f,
		"https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}
	return oauthConfig, nil
}
