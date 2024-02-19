package internal

import (
	"ci6ndex/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
)

type AppConfig struct {
	DiscordToken                   string `mapstructure:"DISCORD_API_TOKEN"`
	DatabaseUrl                    string `mapstructure:"POSTGRES_URL"`
	GoogleCloudCredentialsLocation string `mapstructure:"GCLOUD_CREDS_LOC"`
	CivRankingSheetId              string `mapstructure:"RANKING_SHEET_ID"`
	BotApplicationID               string `mapstructure:"DISCORD_BOT_APPLICATION_ID"`
}

var config AppConfig
var db *DatabaseOperations

const (
	Bot    string = "bot"
	Server string = "server"
)

var server http.Server
var route = mux.NewRouter()

var disc *discordgo.Session

func Start() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("failed to load configuration, error=%w", err))
	}

	db, err = newDBConnection(config.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database, error=%w", err))
	}
	//slog.Info("db start OK")
}

func StartBot() {

	slog.Info("initializing discord bot...")
	disc, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		slog.Error("could not start discord client, exiting", "error", err)
		os.Exit(1)
	}

	disc.Identify.Intents = discordgo.IntentsGuildMessages
	disc.AddHandler(ready)

	err = disc.Open()
	// TODO: These need to be wiped + recreated on startup
	//DeleteDiscordCommands(nil, nil)
	//InitializeDiscordCommands(nil, nil)

	if err != nil {
		slog.Error("could not open connection to discord, exiting", "error", err)
		os.Exit(1)
	}
}

func StartServer() {
	slog.Info("starting http server...")

	route.HandleFunc("/health", Health).Methods("GET")
	route.HandleFunc("/users", CreateUser).Methods("PUT")
	route.HandleFunc("/users/bulk", CreateUsers).Methods("PUT")

	route.HandleFunc("/draft_strategies", GetDraftStrategies).Methods("GET")
	route.HandleFunc("/draft_strategies/{name}", GetDraftStrategy).Methods("GET")

	route.HandleFunc("/drafts", CreateDraft).Methods("PUT")
	route.HandleFunc("/drafts/{draftId}/picks", SubmitDraftPick).Methods("PUT")

	route.HandleFunc("/rankings", RefreshRankings).Methods("POST")

	route.HandleFunc("/discord/commands", GetDiscordCommands).Methods("GET")
	route.HandleFunc("/discord/commands", InitializeDiscordCommands).Methods("POST")
	route.HandleFunc("/discord/commands", DeleteDiscordCommands).Methods("DELETE")

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

func DeleteDiscordCommands(w http.ResponseWriter, req *http.Request) {
	commands, err := disc.ApplicationCommands(config.BotApplicationID, "")
	if err != nil {
		var derr *discordgo.RESTError
		if errors.As(err, &derr) {
			if derr.Response.StatusCode == 404 {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode("could not find commands for guild")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(derr)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		return
	}

	for _, c := range commands {
		err = disc.ApplicationCommandDelete(config.BotApplicationID, "", c.ID)
		if err != nil {
			var derr *discordgo.RESTError
			if errors.As(err, &derr) {
				if derr.Response.StatusCode == 404 {
					w.WriteHeader(http.StatusNotFound)
					_ = json.NewEncoder(w).Encode("could not find commands for guild")
					return
				}
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(derr)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(err)
			}
		}
		slog.Info("removed command", "command", c.Name)
	}
}

func InitializeDiscordCommands(w http.ResponseWriter, req *http.Request) {
	ccmds, err := AttachSlashCommands(disc, &config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errors.Join(errors.New("could not attach slash commands"), err))
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(ccmds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func GetDiscordCommands(w http.ResponseWriter, req *http.Request) {
	commands, err := disc.ApplicationCommands(config.BotApplicationID, "")
	if err != nil {
		var derr *discordgo.RESTError
		if errors.As(err, &derr) {
			if derr.Response.StatusCode == 404 {
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode("could not find commands for guild")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(derr)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(err)
		}
		return
	}

	err = json.NewEncoder(w).Encode(commands)
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

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

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	var users []domain.CreateUsersParams
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	user, err := db.queries.CreateUsers(r.Context(), users)
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

type CreateDraftRequest struct {
	DraftStrategy string `json:"draft_strategy"`
}

func CreateDraft(w http.ResponseWriter, req *http.Request) {
	var cdr CreateDraftRequest
	err := json.NewDecoder(req.Body).Decode(&cdr)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode("could not parse strategy from request body")
		return
	}

	strategy := cdr.DraftStrategy

	_, err = db.queries.GetDraftStrategy(req.Context(), strategy)
	if err != nil {
		w.WriteHeader(422)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("draft_strategy=%v does not exist", strategy))
		return
	}

	draft, err := db.queries.CreateDraft(req.Context(), strategy)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return

	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(draft)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

type Leader struct {
	Name string `json:"name"`
	Civ  string `json:"civ"`
}

// SubmitDraftPickRequest DiscordUser is case sensitive
type SubmitDraftPickRequest struct {
	Leader      Leader   `json:"leader"`
	DiscordUser string   `json:"discord_user"`
	Offered     []Leader `json:"offered"`
}

func SubmitDraftPick(w http.ResponseWriter, r *http.Request) {
	var sdp SubmitDraftPickRequest
	err := json.NewDecoder(r.Body).Decode(&sdp)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode("could not parse draft pick from request body")
		return
	}

	vars := mux.Vars(r)
	draftId, err := strconv.ParseInt(vars["draftId"], 10, 64)
	if err != nil {
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode("could not parse draft id from request path")
		return
	}

	_, err = db.queries.GetDraft(r.Context(), draftId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(422)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("draft_id=%v does not exist", draftId))
			return
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	user, err := db.queries.GetUserByDiscordName(r.Context(), sdp.DiscordUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(422)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("discord_user=%v does not exist", sdp.DiscordUser))
			return
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	leader, err := db.queries.GetLeaderByNameAndCiv(r.Context(), domain.GetLeaderByNameAndCivParams{
		LeaderName: strings.ToUpper(sdp.Leader.Name),
		CivName:    strings.ToUpper(sdp.Leader.Civ),
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(422)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("leader=%v does not exist", sdp.Leader))
			return
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	offered, err := json.Marshal(sdp.Offered)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("could not serialize offered leaders to database"))
		return
	}

	pick, err := db.queries.SubmitDraftPick(r.Context(), domain.SubmitDraftPickParams{
		DraftID:  draftId,
		LeaderID: pgtype.Int8{Int64: leader.ID},
		UserID:   user.ID,
		Offered:  offered,
	})

	if err != nil {
		var sqlErr *pgconn.PgError
		if errors.As(err, &sqlErr) {
			if sqlErr.Code == "23505" {
				w.WriteHeader(409)
				err = json.NewEncoder(w).Encode(fmt.Sprint("user has already submitted a pick for this draft"))
				return
			}
		} else {
			w.WriteHeader(500)
			_ = json.NewEncoder(w).Encode(err)
			return
		}
	}

	w.WriteHeader(201)
	err = json.NewEncoder(w).Encode(pick)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(fmt.Sprint("failed to encode draft pick, but it was submitted successfully"))
	}

}

func GetDraftStrategies(w http.ResponseWriter, r *http.Request) {
	strats, err := db.queries.GetDraftStrategies(r.Context())
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(200)
	if strats == nil {
		strats = []domain.Ci6ndexDraftStrategy{}
	}
	err = json.NewEncoder(w).Encode(strats)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

func GetDraftStrategy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// We can defer sql sanitization to pgx here because they use prepared statements
	// and our sqlc generated code "parameterizes our parameters" (i.e. it uses $1, $2, etc.)
	// See more:
	// * https://github.com/jackc/pgx/wiki/Automatic-Prepared-Statement-Caching#automatic-prepared-statement-caching
	strat, err := db.queries.GetDraftStrategy(r.Context(), vars["name"])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(404)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("draft_strategy=%v not found", vars["name"]))
			return
		}
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(strat)
	if err != nil {
		w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
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

func newDBConnection(dbUrl string) (*DatabaseOperations, error) {
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

func (s DatabaseOperations) Health() error {
	err := s.db.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s DatabaseOperations) Close() {
	s.db.Close()
}

func (r Ranking) ToRankingDBParam(ctx context.Context) (domain.CreateRankingParams, error) {
	user, err := db.queries.GetUserByName(ctx, r.Player)
	if err != nil {
		return domain.CreateRankingParams{}, errors.New(fmt.Sprintf("could not find user=%v from google sheets in local database", r.Player))
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

	l, err := db.queries.GetLeaderByNameAndCiv(ctx, domain.GetLeaderByNameAndCivParams{
		LeaderName: strings.ToUpper(leader),
		CivName:    strings.ToUpper(civ),
	})

	if err != nil {
		return domain.CreateRankingParams{}, err
	}

	return domain.CreateRankingParams{
		UserID:   user.ID,
		Tier:     r.Tier,
		LeaderID: l.ID,
	}, nil
}
