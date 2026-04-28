package ci6ndex

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"ci6ndex/ci6ndex/generated"

	_ "github.com/mattn/go-sqlite3"
	goose "github.com/pressly/goose/v3"
)

const testGuildID uint64 = 999999
const testMemoryDSN = "file:test?mode=memory&cache=shared"

var (
	testDB *DB
	testC  *Ci6ndex
)

func TestMain(m *testing.M) {
	writeConn, err := sql.Open("sqlite3", testMemoryDSN)
	if err != nil {
		slog.Error("failed to open test write connection", "error", err)
		os.Exit(1)
	}
	defer writeConn.Close()

	readConn, err := sql.Open("sqlite3", testMemoryDSN)
	if err != nil {
		slog.Error("failed to open test read connection", "error", err)
		os.Exit(1)
	}
	defer readConn.Close()

	// Configure goose to read migrations from project root
	goose.SetBaseFS(os.DirFS(".."))
	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("failed to set goose dialect", "error", err)
		os.Exit(1)
	}

	if err := goose.Up(writeConn, "sql/migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	// Enforce single writer (matches production)
	writeConn.SetMaxOpenConns(1)

	queries := generated.New(readConn)
	writes := generated.New(writeConn)

	testDB = &DB{
		readConn:  readConn,
		writeConn: writeConn,
		Queries:   queries,
		Writes:    writes,
	}

	// Seed test data
	if err := seedTestData(); err != nil {
		slog.Error("failed to seed test data", "error", err)
		os.Exit(1)
	}

	testC = &Ci6ndex{
		Connections: map[uint64]*DB{testGuildID: testDB},
		Path:        "",
	}

	code := m.Run()
	os.Exit(code)
}

// seedTestData inserts players, draft, and draft_registry rows.
// 18 players are registered so exhaustion tests can be performed.
func seedTestData() error {
	ctx := context.Background()

	// 18 players allows exhaustion tests (85 eligible leaders / 5 per player = 17 max)
	players := make([]generated.AddPlayerParams, 18)
	for i := range players {
		players[i] = generated.AddPlayerParams{
			ID:            int64(1000 + i),
			Username:      "test_player_" + string(rune('0'+i)),
			GlobalName:    sql.NullString{},
			DiscordAvatar: sql.NullString{},
		}
	}

	for _, p := range players {
		if err := testDB.Writes.AddPlayer(ctx, p); err != nil {
			return err
		}
	}

	draft, err := testDB.Writes.CreateActiveDraft(ctx)
	if err != nil {
		return err
	}

	for _, p := range players {
		if _, err := testDB.Writes.AddPlayerToDraft(ctx, generated.AddPlayerToDraftParams{
			DraftID:  draft.ID,
			PlayerID: p.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}

// standardRules returns the 5-rule set used in production rolls.
func standardRules() []Rule {
	return []Rule{
		&MinTierRule{MinTier: 3},
		&NoOpRule{},
		&NoOpRule{},
		&NoOpRule{},
		&NoOpRule{},
	}
}

func TestRollForPlayers_Basic(t *testing.T) {
	ctx := context.Background()
	players, err := testDB.Queries.GetPlayersFromActiveDraft(ctx)
	if err != nil {
		t.Fatalf("failed to get players: %v", err)
	}

	playerIds := make([]int64, 3)
	for i := range 3 {
		playerIds[i] = players[i].ID
	}

	offerings, err := testC.RollForPlayers(testGuildID, playerIds, standardRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(offerings) != len(playerIds) {
		t.Fatalf("expected %d offerings, got %d", len(playerIds), len(offerings))
	}

	seen := make(map[int64]bool)
	for _, o := range offerings {
		if len(o.Leaders) != 5 {
			t.Fatalf("expected 5 leaders per offering, got %d", len(o.Leaders))
		}
		for _, l := range o.Leaders {
			if seen[l.ID] {
				t.Fatalf("leader %d (%s) assigned to multiple players", l.ID, l.LeaderName)
			}
			seen[l.ID] = true
		}
	}
}

func TestRollForPlayers_MinTier(t *testing.T) {
	ctx := context.Background()
	players, err := testDB.Queries.GetPlayersFromActiveDraft(ctx)
	if err != nil {
		t.Fatalf("failed to get players: %v", err)
	}

	playerIds := make([]int64, 3)
	for i := range 3 {
		playerIds[i] = players[i].ID
	}

	offerings, err := testC.RollForPlayers(testGuildID, playerIds, standardRules())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, o := range offerings {
		hasLowTier := false
		for _, l := range o.Leaders {
			if l.Tier <= 3 {
				hasLowTier = true
				break
			}
		}
		if !hasLowTier {
			t.Fatalf("player %d has no leader with tier <= 3", o.Player.ID)
		}
	}
}

func TestRollForPlayers_RanOutOfChoices(t *testing.T) {
	ctx := context.Background()
	players, err := testDB.Queries.GetPlayersFromActiveDraft(ctx)
	if err != nil {
		t.Fatalf("failed to get players: %v", err)
	}

	// Use only NoOpRules to avoid MinTier constraint exhaustion
	rules := []Rule{
		&NoOpRule{},
		&NoOpRule{},
		&NoOpRule{},
		&NoOpRule{},
		&NoOpRule{},
	}

	playerIds := make([]int64, len(players))
	for i, p := range players {
		playerIds[i] = p.ID
	}

	_, err = testC.RollForPlayers(testGuildID, playerIds, rules)
	if err == nil {
		t.Fatal("expected RanOutOfChoicesError with too many players, got nil")
	}

	if _, ok := err.(RanOutOfChoicesError); !ok {
		t.Fatalf("expected RanOutOfChoicesError, got %T: %v", err, err)
	}
}
