package pkg

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func AddUsersFromFile(path string, db *DatabaseOperations) error {
	exists := fileExists(path)
	if !exists {
		return fmt.Errorf("file does not exist")
	}
	if !strings.HasSuffix(path, ".json") {
		return fmt.Errorf("file must be a json file")
	}

	users, err := parseUsers(path)
	if err != nil {
		return err
	}

	dbUsers, err := db.Queries.CreateUsers(context.Background(), users)

	if err != nil {
		return err
	}
	slog.Info("Added users", "count", dbUsers)
	return nil
}

func fileExists(relPath string) bool {
	_, err := os.Stat(relPath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		slog.Error("Error checking if file exists", "err", err)
		return false
	}
	return true
}

func parseUsers(path string) ([]domain.CreateUsersParams, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []domain.CreateUsersParams
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
