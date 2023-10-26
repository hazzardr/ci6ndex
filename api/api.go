package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

type AppConfig struct {
	DiscordToken                   string `mapstructure:"DISCORD_API_TOKEN"`
	DatabaseUrl                    string `mapstructure:"POSTGRES_URL"`
	GoogleCloudCredentialsLocation string `mapstructure:"GCLOUD_CREDS_LOC"`
	CivRankingSheetId              string `mapstructure:"RANKING_SHEET_ID"`
}

var config AppConfig
var storage *Storage

type Mode string

const (
	Bot    Mode = "bot"
	Server Mode = "server"
)

var server http.Server
var r = mux.NewRouter()

var d discordgo.Session

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

	//TODO:
	storage, err = NewDBConnection(config.DatabaseUrl)
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
	d.AddHandler(messageCreate)

	err = d.Open()
	if err != nil {
		slog.Error("could not open connection to discord, exiting", "error", err)
		os.Exit(1)
	}
}

func StartServer() {
	slog.Info("starting http server...")
	server.Handler = r
	server.Addr = ":8080"

	r.HandleFunc("/health", Health).Methods("GET")
	r.HandleFunc("/rankings", GetRankings).Methods("GET")
	http.Handle("/", r)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-stop

	StopServer(0)
}

func GetRankings(w http.ResponseWriter, r *http.Request) {
	err := updateLocalRankings(&config)
	if err != nil {
		slog.Error("could not get rankings", "error", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func StopServer(code int) {
	slog.Info("shutting down application")
	err := d.Close()
	if err != nil {
		slog.Warn("did not shut down application gracefully", "error", err)
	} else {
		slog.Info("discord connection shut down successfully")
	}
	err = server.Shutdown(context.Background())
	if err != nil {
		slog.Warn("did not shut down application gracefully", "error", err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {
	err := storage.Health()
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

type Storage struct {
	db *pgxpool.Pool
}

func (s Storage) Health() error {
	err := s.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func NewDBConnection(dbUrl string) (*Storage, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Storage{db: conn}, nil
}
