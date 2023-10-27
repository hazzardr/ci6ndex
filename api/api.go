package api

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
)

type AppConfig struct {
	DiscordToken                   string `mapstructure:"DISCORD_API_TOKEN"`
	DatabaseUrl                    string `mapstructure:"POSTGRES_URL"`
	GoogleCloudCredentialsLocation string `mapstructure:"GCLOUD_CREDS_LOC"`
	CivRankingSheetId              string `mapstructure:"RANKING_SHEET_ID"`
}

var config AppConfig
var db *DatabaseOperations

type Mode string

const (
	Bot    Mode = "bot"
	Server Mode = "server"
)

var server http.Server
var route = mux.NewRouter()

var disc discordgo.Session

func Start(mode string) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	db, err = NewDBConnection(config.DatabaseUrl)
	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	switch mode {
	case string(Server):
		slog.Info("starting in server mode only")
		StartServer()
	case string(Bot):
		StartBot()
		StartServer()
	default:
		slog.Error("unsupported mode passed. exiting", "mode", mode)
		os.Exit(1)
	}
}

func StartBot() {
	slog.Info("initializing discord bot...")
	d, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		slog.Error("could not start discord client, exiting", "error", err)
		os.Exit(1)
	}

	d.Identify.Intents = discordgo.IntentsGuildMessages
	d.AddHandler(ready)
	//disc.AddHandler(messageCreate)

	err = d.Open()
	if err != nil {
		slog.Error("could not open connection to discord, exiting", "error", err)
		os.Exit(1)
	}
}

func StartServer() {
	slog.Info("starting http server...")

	route.HandleFunc("/health", Health).Methods("GET")
	route.HandleFunc("/users", CreateUser).Methods("PUT")

	//r.HandleFunc("/rankings", GetRankings).Methods("GET")
	route.HandleFunc("/rankings", RefreshRankings).Methods("POST")

	server = http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-stop

	StopServer(0)
}

func Health(w http.ResponseWriter, r *http.Request) {
	err := db.Health()
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode("OK")
	if err != nil {
		w.WriteHeader(500)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var params domain.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	user, err := db.queries.CreateUser(r.Context(), params)
	if err != nil {
		var sqlErr *pgconn.PgError
		if errors.As(err, &sqlErr) {
			if sqlErr.Code == "23505" {
				w.WriteHeader(409)
				_ = json.NewEncoder(w).Encode("user already exists")
				return
			}
		} else {
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
		return
	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

func GetRankings(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func RefreshRankings(w http.ResponseWriter, req *http.Request) {
	ranks, err := getRankingsFromSheets(&config, req.Context())
	if err != nil {
		slog.Error("could not get rankings", "error", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	err = db.queries.DeleteRankings(req.Context())
	if err != nil {
		slog.Error("could not delete existing rankings", "error", err)
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	for _, r := range ranks {
		dbp, err := r.ToRankingDBParam(req.Context())
		if err != nil {
			slog.Error("could not convert gsheet ranking to db param", "error", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(err)
			return
		}

		_, err = db.queries.CreateRanking(req.Context(), dbp)
		if err != nil {
			slog.Error("could not create ranking", "error", err)
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
	}

	w.WriteHeader(200)
	success := "successfully refreshed rankings from google sheets"
	slog.Info(success, "ranks_added", len(ranks))
	err = json.NewEncoder(w).Encode(success)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

func StopServer(code int) {
	slog.Info("shutting down application")
	err := disc.Close()
	if err != nil {
		slog.Warn("did not shut down application gracefully", "error", err)
	} else {
		slog.Info("discord connection shut down successfully")
	}

	slog.Info("closing database connection")
	db.Close()
	slog.Info("database connection closed successfully")

	err = server.Shutdown(context.Background())
	if err != nil {
		slog.Warn("did not shut down application gracefully", "error", err)
	}
}

type DatabaseOperations struct {
	db      *pgxpool.Pool
	queries *domain.Queries
}

func (s DatabaseOperations) Health() error {
	err := s.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func NewDBConnection(dbUrl string) (*DatabaseOperations, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	q := domain.New(conn)

	return &DatabaseOperations{db: conn, queries: q}, nil
}

func (s DatabaseOperations) Close() {
	s.db.Close()
}

func (r Ranking) ToRankingDBParam(ctx context.Context) (domain.CreateRankingParams, error) {
	user, err := db.queries.GetUserByName(ctx, r.Player)
	if err != nil {
		return domain.CreateRankingParams{}, err
	}

	re, err := regexp.Compile(`^(.*?) \((.*?)\)$`)
	if err != nil {
		return domain.CreateRankingParams{}, err
	}
	matches := re.FindStringSubmatch(r.CombinedLeaderCiv)

	var civ string
	var leader string
	if len(matches) == 3 {
		civ = matches[1]
		leader = matches[2]
	} else {
		return domain.CreateRankingParams{}, errors.New("could not parse civ and leader from google sheets cell")
	}

	l, err := db.queries.GetLeader(ctx, domain.GetLeaderParams{
		LeaderName: leader,
		CivName:    civ,
	})

	if err != nil {
		return domain.CreateRankingParams{}, err
	}

	return domain.CreateRankingParams{
		UserID:   int32(user.ID),
		Tier:     r.Tier,
		LeaderID: int32(l.ID),
	}, nil
}
