package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
)

type Config struct {
	DiscordToken string `mapstructure:"DISCORD_API_TOKEN"`
}

var config Config

type Mode string

const (
	Bot    Mode = "bot"
	Server Mode = "server"
)

var server http.Server
var r = mux.NewRouter()

var d discordgo.Session

func Start(mode string) {
	//TODO: get discord token / register app
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
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

	r.HandleFunc("/health", Health).Methods("GET")

	http.Handle("/", r)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	<-stop

	StopServer(0)
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
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("OK")
}
