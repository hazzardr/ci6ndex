package discord

import (
	"ci6ndex/domain"
	"ci6ndex/internal"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"os"
	"os/signal"
)

const (
	RollCivs = "roll-civs"
)

type Bot struct {
	s      *discordgo.Session
	db     *internal.DatabaseOperations
	config *internal.AppConfig
}

type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func NewBot(db *internal.DatabaseOperations, config *internal.AppConfig) (*Bot, error) {
	s, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, fmt.Errorf("could not start discord client: %w", err)
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages
	return &Bot{
		s:      s,
		db:     db,
		config: config,
	}, nil
}

func (bot *Bot) Start() error {
	bot.s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info(fmt.Sprintf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})
	err := bot.s.Open()
	if err != nil {
		return fmt.Errorf("cannot open the session: %v", err)
	}

	defer bot.s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	slog.Info("bot initialized and ready to receive events")
	<-stop
	slog.Info("received interrupt signal, shutting down")
	return nil
}

// RegisterSlashCommands attaches all slash commands to the bot. Database has to be initialized first.
func (bot *Bot) RegisterSlashCommands(guild string) ([]*discordgo.ApplicationCommand, error) {
	err := bot.db.Health()
	if err != nil {
		return nil, fmt.Errorf("can't attach commands prior to db being initialized: %w", err)
	}

	commands := getDraftCommands()
	handlers := getDraftHandlers(bot.db, bot.config)

	for _, c := range commands {
		_, err := bot.s.ApplicationCommandCreate(bot.config.BotApplicationID, guild, c)
		if err != nil {
			slog.Error("could not create (/) command", "command", c.Name, "error", err)
			return nil, err
		}
		slog.Info("registered", "command", c.Name)
	}
	slog.Info("all (/) commands attached")
	bot.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	return commands, nil
}

func (bot *Bot) RemoveSlashCommands(guild string) error {
	commands, err := bot.s.ApplicationCommands(bot.config.BotApplicationID, guild)
	if err != nil {
		return err
	}
	if nil == commands || len(commands) == 0 {
		slog.Info("no commands to remove", "guildId", guild)
		return nil
	}
	for _, c := range commands {
		err = bot.s.ApplicationCommandDelete(bot.config.BotApplicationID, guild, c.ID)
		if err != nil {
			return err
		}
		slog.Info("removed command", "command", c.Name, "guildId", guild)
	}

	return nil
}

func (bot *Bot) rollCivs(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
	ctx := context.Background()
	drafts, err := bot.db.Queries.GetActiveDrafts(ctx)
	if err != nil {
		bot.reportError("error checking active drafts", err, i)
		return
	}

	var activeDraft domain.Ci6ndexDraft

	if len(drafts) == 0 {
		_, err = bot.s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "There is no active draft. Roll will not be attached to a game.",
		})
		if err != nil {
			slog.Error("error responding to user", "error", err)
		}

		// dummy draft as a default
		activeDraft, err = bot.db.Queries.GetDraft(ctx, -1)
		if err != nil {
			bot.reportError("error fetching dummy draft", err, i)
			return
		}

	}

	if len(drafts) > 1 {
		bot.reportError("There are multiple active drafts. This should not be possible", nil, i)
		return
	}

	if len(drafts) == 1 {
		_, err = bot.s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Rolled civs will be attached to the active draft.",
		})

		if err != nil {
			slog.Error("error responding to user", "error", err)
		}

		activeDraft = drafts[0]
	}
	leaders, err := bot.db.Queries.GetLeaders(ctx)
	if err != nil {
		bot.reportError("error fetching leaders", err, i)
		return
	}

	strat, err := bot.db.Queries.GetDraftStrategy(ctx, activeDraft.DraftStrategy)
	shuffler := internal.NewCivShuffler(leaders, activeDraft.Players, strat, bot.db)
	offers, err := shuffler.Shuffle()
	if err != nil {
		bot.reportError("error shuffling civs", err, i)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("The following picks were rolled: %v", offers),
		},
	})

	if err != nil {
		slog.Error("error responding to user", "error", err)
	}
}

//	func submitPicks(s *discordgo.Session, i *discordgo.InteractionCreate) {
//		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: "this is a pick field",
//				Components: []discordgo.MessageComponent{
//					discordgo.SelectMenu{
//						MenuType:  discordgo.SelectMenuType(discordgo.SelectMenuComponent),
//						MaxValues: 1,
//						Disabled:  false,
//						Options: []discordgo.SelectMenuOption{
//							{
//								Label:       "test1",
//								Value:       "test1 val",
//								Description: "test1 desc",
//							},
//						},
//					},
//				},
//			},
//		})
//
//		if err != nil {
//			ReportError("error picking civs", err, s, i)
//		}
//	}
//
//	func players(s *discordgo.Session, i *discordgo.InteractionCreate) {
//		slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//		users, err := db.Queries.GetUsers(context.Background())
//		if err != nil {
//			slog.Error("error getting players", "error", err)
//			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: fmt.Sprintf("Error getting players: %s", err.Error()),
//				},
//			})
//			if err != nil {
//				slog.Error("error responding to user", "error", err)
//			}
//			return
//		}
//		var playerNames []string
//		for _, p := range users {
//			playerNames = append(playerNames, p.Name)
//		}
//		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: fmt.Sprintf("Players eligible for draft: %s", strings.Join(playerNames, ", ")),
//			},
//		})
//		if err != nil {
//			slog.Error(err.Error())
//		}
//	}
//
//	func startDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
//		slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//
//		options := i.ApplicationCommandData().Options
//		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
//
//		for _, opt := range options {
//			optionMap[opt.Name] = opt
//		}
//
//		strat := optionMap["draft-strategy"].StringValue()
//
//		if strat == "" {
//			slog.Error("no strategy provided for draft - this should not be possible")
//			return
//		}
//
//		ds, err := db.Queries.GetDraftStrategy(context.Background(), strat)
//
//		if err != nil {
//			ReportError("error fetching draft strategy", err, s, i)
//			return
//		}
//
//		actives, err := db.Queries.GetActiveDrafts(context.Background())
//
//		if err != nil {
//			ReportError("error fetching active draft", err, s, i)
//			return
//		}
//
//		if len(actives) > 0 {
//			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//				Type: discordgo.InteractionResponseChannelMessageWithSource,
//				Data: &discordgo.InteractionResponseData{
//					Content: "There is already an active draft. Please end it before starting a new one.",
//				},
//			})
//			if err != nil {
//				slog.Error("error responding to user", "error", err)
//			}
//			return
//		}
//
//		draft, err := db.Queries.CreateDraft(context.Background(), ds.Name)
//
//		if err != nil {
//			ReportError("error creating draft", err, s, i)
//			return
//		}
//
//		slog.Info("draft created", "draft", draft.ID, "strategy", ds.Name, "startedBy", i.Interaction.Member.User.Username)
//		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: fmt.Sprintf("Draft #%v %s started by user %s. %s", draft.ID, ds.Name, i.Interaction.Member.User.Username, ds.Description),
//			},
//		})
//		if err != nil {
//			slog.Error(err.Error())
//		}
//	}
func basicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ci6ndex (Civ VI Index) is a bot for managing civ 6 drafts. Use /start-draft to start a draft, or /roll-civs to assign civs to players",
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

//	func ready(s *discordgo.Session, e *discordgo.Ready) {
//		err := RemoveSlashCommands()
//		if err != nil {
//			slog.Error("could not remove slash commands", "error", err)
//			os.Exit(1)
//		}
//		_, err = RegisterSlashCommands(s)
//		if err != nil {
//			slog.Error("could not attach slash commands", "error", err)
//			os.Exit(1)
//		}
//
//		err = s.UpdateGameStatus(0, "/ci6ndex")
//		if err != nil {
//			slog.Warn("could not update discord status on startup")
//		}
//		slog.Info("bot initialized and ready to receive events")
//	}
func (bot *Bot) reportError(msg string, err error, i *discordgo.InteractionCreate) {
	slog.Error(msg, "error", err)
	_, err = bot.s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Something went wrong",
	})
	if err != nil {
		slog.Error("error responding to user", "error", err)
	}
}

//
////
////func StartBot() {
////
////	slog.Info("initializing discord bot...")
////	disc, err := discordgo.New("Bot " + config.DiscordToken)
////	if err != nil {
////		slog.Error("could not start discord client, exiting", "error", err)
////		os.Exit(1)
////	}
////
////	disc.Identify.Intents = discordgo.IntentsGuildMessages
////	disc.AddHandler(ready)
////
////	err = disc.Open()
////
////	if err != nil {
////		slog.Error("could not open connection to discord, exiting", "error", err)
////		os.Exit(1)
////	}
////}
////
////func DeleteDiscordCommands(w http.ResponseWriter, req *http.Request) {
////	err := RemoveSlashCommands()
////	if err != nil {
////		var derr *discordgo.RESTError
////		if errors.As(err, &derr) {
////			if derr.Response.StatusCode == 404 {
////				w.WriteHeader(http.StatusNotFound)
////				_ = json.NewEncoder(w).Encode("could not find commands for guild")
////				return
////			}
////			w.WriteHeader(http.StatusInternalServerError)
////			_ = json.NewEncoder(w).Encode(derr)
////		} else {
////			w.WriteHeader(http.StatusInternalServerError)
////			_ = json.NewEncoder(w).Encode(err)
////		}
////		return
////	}
////
////	w.WriteHeader(http.StatusOK)
////	err = json.NewEncoder(w).Encode("successfully deleted commands")
////}
////
////func InitializeDiscordCommands(w http.ResponseWriter, req *http.Request) {
////	ccmds, err := RegisterSlashCommands(disc)
////	if err != nil {
////		w.WriteHeader(http.StatusInternalServerError)
////		_ = json.NewEncoder(w).Encode(errors.Join(errors.New("could not attach slash commands"), err))
////	}
////	w.WriteHeader(http.StatusOK)
////	err = json.NewEncoder(w).Encode(ccmds)
////	if err != nil {
////		w.WriteHeader(http.StatusInternalServerError)
////	}
////
////}
////
////func GetDiscordCommands(w http.ResponseWriter, req *http.Request) {
////	commands, err := disc.ApplicationCommands(config.BotApplicationID, "")
////	if err != nil {
////		var derr *discordgo.RESTError
////		if errors.As(err, &derr) {
////			if derr.Response.StatusCode == 404 {
////				w.WriteHeader(http.StatusNotFound)
////				_ = json.NewEncoder(w).Encode("could not find commands for guild")
////				return
////			}
////			w.WriteHeader(http.StatusInternalServerError)
////			_ = json.NewEncoder(w).Encode(derr)
////		} else {
////			w.WriteHeader(http.StatusInternalServerError)
////			_ = json.NewEncoder(w).Encode(err)
////		}
////		return
////	}
////
////	err = json.NewEncoder(w).Encode(commands)
////	w.WriteHeader(http.StatusOK)
////	if err != nil {
////		w.WriteHeader(http.StatusInternalServerError)
////		_ = json.NewEncoder(w).Encode(err)
////		return
////	}
////
////}
