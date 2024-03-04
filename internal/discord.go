package internal

//
//import (
//	"ci6ndex/domain"
//	"context"
//	"fmt"
//	"github.com/bwmarrin/discordgo"
//	"log/slog"
//	"os"
//	"strings"
//)
//
//const (
//	StartDraft = "start-draft"
//	Players    = "players"
//	RollCivs   = "roll-civs"
//	SubmitPick = "submit-pick"
//)
//
//type DiscordBot struct {
//	session *discordgo.Session
//	db *DatabaseOperations
//}
//
//// AttachSlashCommands attaches all slash commands to the bot. Database has to be initialized first.
//func AttachSlashCommands(s *discordgo.Session, db *DatabaseOperations) ([]*discordgo.ApplicationCommand, error) {
//	err := db.Health()
//	if err != nil {
//		return nil, fmt.Errorf("can't attach commands prior to db being initialized: %w", err)
//	}
//	strategies, err := db.Queries.GetDraftStrategies(context.Background())
//	if err != nil {
//		return nil, fmt.Errorf("couldn't get required defaults: %w", err)
//	}
//
//	var stratOptions []*discordgo.ApplicationCommandOptionChoice
//	for _, strategy := range strategies {
//		stratOptions = append(stratOptions, &discordgo.ApplicationCommandOptionChoice{
//			Name:  strategy.Name,
//			Value: strategy.Name,
//		})
//	}
//
//	draftStrategyOptions := []*discordgo.ApplicationCommandOption{
//		{
//			Type:        discordgo.ApplicationCommandOptionString,
//			Name:        "draft-strategy",
//			Description: "How to decide which civs to roll",
//			Choices:     stratOptions,
//			Required:    true,
//		},
//	}
//
//	commands := []*discordgo.ApplicationCommand{
//		{
//			Name:        StartDraft,
//			Description: "Initialize a draft. You can then roll civs for players and they can submit their picks.",
//			Options:     draftStrategyOptions,
//		},
//		{
//			Name:        RollCivs,
//			Description: "Roll civs for all players",
//			Options:     draftStrategyOptions,
//		},
//		{
//			Name:        Players,
//			Description: "Get a list of players eligible for a draft",
//		},
//		{
//			Name:        "ci6ndex",
//			Description: "Get information about the bot",
//		},
//		{
//			Name:        SubmitPick,
//			Description: "Submit a pick for a draft",
//		},
//	}
//	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
//		"ci6ndex":  basicCommand,
//		StartDraft: startDraft,
//		Players:    players,
//		RollCivs:   rollCivs,
//		SubmitPick: submitPicks,
//	}
//
//	var ccommands []*discordgo.ApplicationCommand
//	for _, c := range commands {
//		ccmd, err := s.ApplicationCommandCreate(config.BotApplicationID, "", c)
//		if err != nil {
//			slog.Error("could not create slash command", "command", c.Name, "error", err)
//		}
//		slog.Info("registered", "command", c.Name)
//		ccommands = append(ccommands, ccmd)
//	}
//	slog.Info("slash commands attached")
//	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
//		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
//			h(s, i)
//		}
//	})
//	return ccommands, nil
//}
//
//func RemoveSlashCommands(disc *discordgo.Session, config AppConfig) error {
//	commands, err := disc.ApplicationCommands(config.BotApplicationID, "")
//	if err != nil {
//		return err
//	}
//	for _, c := range commands {
//		err = disc.ApplicationCommandDelete(config.BotApplicationID, "", c.ID)
//		if err != nil {
//			return err
//		}
//		slog.Info("removed command", "command", c.Name)
//	}
//
//	return nil
//}
//
//// TODO: switch to followups so we can send multiple messages
//func (bot *DiscordBot) rollCivs(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//	drafts, err := bot.db.Queries.GetActiveDrafts(context.Background())
//	if err != nil {
//		ReportError("error checking active drafts", err, s, i)
//		return
//	}
//
//	var activeDraft domain.Ci6ndexDraft
//
//	if len(drafts) == 0 {
//		//err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		//	Type: discordgo.InteractionResponseChannelMessageWithSource,
//		//	Data: &discordgo.InteractionResponseData{
//		//		Content: "There is no active draft. These results will not be attached to a game",
//		//	},
//		//})
//		//if err != nil {
//		//	slog.Error("error responding to user", "error", err)
//		//}
//
//		// dummy draft as a default
//		activeDraft = domain.Ci6ndexDraft{
//			ID:            -1,
//			DraftStrategy: "RandomPickPool3Standard",
//			Active:        true,
//		}
//	}
//
//	if len(drafts) > 1 {
//		ReportError("there are multiple active drafts. This should not be possible", nil, s, i)
//		return
//	}
//
//	if len(drafts) == 1 {
//		//err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		//	Type: discordgo.InteractionResponseChannelMessageWithSource,
//		//	Data: &discordgo.InteractionResponseData{
//		//		Content: "Rolled civs will be attached to the active draft",
//		//	},
//		//})
//
//		if err != nil {
//			slog.Error("error responding to user", "error", err)
//		}
//
//		activeDraft = drafts[0]
//	}
//
//	picks, err := OfferPicks(bot.db, activeDraft, 3)
//	if err != nil {
//		ReportError("error rolling civs", err, s, i)
//		return
//	}
//
//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: fmt.Sprintf("The following picks were rolled: %v", picks),
//		},
//	})
//
//	if err != nil {
//		slog.Error("error responding to user", "error", err)
//	}
//}
//
//func submitPicks(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: "this is a pick field",
//			Components: []discordgo.MessageComponent{
//				discordgo.SelectMenu{
//					MenuType:  discordgo.SelectMenuType(discordgo.SelectMenuComponent),
//					MaxValues: 1,
//					Disabled:  false,
//					Options: []discordgo.SelectMenuOption{
//						{
//							Label:       "test1",
//							Value:       "test1 val",
//							Description: "test1 desc",
//						},
//					},
//				},
//			},
//		},
//	})
//
//	if err != nil {
//		ReportError("error picking civs", err, s, i)
//	}
//}
//
//func players(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//	users, err := db.Queries.GetUsers(context.Background())
//	if err != nil {
//		slog.Error("error getting players", "error", err)
//		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: fmt.Sprintf("Error getting players: %s", err.Error()),
//			},
//		})
//		if err != nil {
//			slog.Error("error responding to user", "error", err)
//		}
//		return
//	}
//	var playerNames []string
//	for _, p := range users {
//		playerNames = append(playerNames, p.Name)
//	}
//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: fmt.Sprintf("Players eligible for draft: %s", strings.Join(playerNames, ", ")),
//		},
//	})
//	if err != nil {
//		slog.Error(err.Error())
//	}
//}
//
//func startDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//
//	options := i.ApplicationCommandData().Options
//	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
//
//	for _, opt := range options {
//		optionMap[opt.Name] = opt
//	}
//
//	strat := optionMap["draft-strategy"].StringValue()
//
//	if strat == "" {
//		slog.Error("no strategy provided for draft - this should not be possible")
//		return
//	}
//
//	ds, err := db.Queries.GetDraftStrategy(context.Background(), strat)
//
//	if err != nil {
//		ReportError("error fetching draft strategy", err, s, i)
//		return
//	}
//
//	actives, err := db.Queries.GetActiveDrafts(context.Background())
//
//	if err != nil {
//		ReportError("error fetching active draft", err, s, i)
//		return
//	}
//
//	if len(actives) > 0 {
//		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//			Type: discordgo.InteractionResponseChannelMessageWithSource,
//			Data: &discordgo.InteractionResponseData{
//				Content: "There is already an active draft. Please end it before starting a new one.",
//			},
//		})
//		if err != nil {
//			slog.Error("error responding to user", "error", err)
//		}
//		return
//	}
//
//	draft, err := db.Queries.CreateDraft(context.Background(), ds.Name)
//
//	if err != nil {
//		ReportError("error creating draft", err, s, i)
//		return
//	}
//
//	slog.Info("draft created", "draft", draft.ID, "strategy", ds.Name, "startedBy", i.Interaction.Member.User.Username)
//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: fmt.Sprintf("Draft #%v %s started by user %s. %s", draft.ID, ds.Name, i.Interaction.Member.User.Username, ds.Description),
//		},
//	})
//	if err != nil {
//		slog.Error(err.Error())
//	}
//}
//
//func basicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
//	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
//	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: "Ci6ndex (Civ VI Index) is a bot for managing civ 6 drafts. Use /start-draft to start a draft, or /roll-civs to assign civs to players",
//		},
//	})
//	if err != nil {
//		slog.Error(err.Error())
//	}
//}
//
//func ready(s *discordgo.Session, e *discordgo.Ready) {
//	err := RemoveSlashCommands()
//	if err != nil {
//		slog.Error("could not remove slash commands", "error", err)
//		os.Exit(1)
//	}
//	_, err = AttachSlashCommands(s)
//	if err != nil {
//		slog.Error("could not attach slash commands", "error", err)
//		os.Exit(1)
//	}
//
//	err = s.UpdateGameStatus(0, "/ci6ndex")
//	if err != nil {
//		slog.Warn("could not update discord status on startup")
//	}
//	slog.Info("bot initialized and ready to receive events")
//}
//
//func ReportError(msg string, err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
//	slog.Error(msg, "error", err)
//	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
//		Type: discordgo.InteractionResponseChannelMessageWithSource,
//		Data: &discordgo.InteractionResponseData{
//			Content: msg,
//		},
//	})
//	if err != nil {
//		slog.Error("error responding to user", "error", err)
//	}
//}
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
////	ccmds, err := AttachSlashCommands(disc)
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
