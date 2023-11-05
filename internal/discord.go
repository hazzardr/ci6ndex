package internal

import (
	"ci6ndex/domain"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"strings"
)

const (
	StartDraft = "start-draft"
	Players    = "players"
	RollCivs   = "roll-civs"
)

// AttachSlashCommands attaches all slash commands to the bot. Database has to be initialized first.
func AttachSlashCommands(s *discordgo.Session, config *AppConfig) ([]*discordgo.ApplicationCommand, error) {
	err := db.Health()
	if err != nil {
		return nil, fmt.Errorf("can't attach commands prior to db being initialized: %w", err)
	}
	strategies, err := db.queries.GetDraftStrategies(context.Background())
	if err != nil {
		return nil, fmt.Errorf("couldn't get required defaults: %w", err)
	}

	var stratOptions []*discordgo.ApplicationCommandOptionChoice
	for _, strategy := range strategies {
		stratOptions = append(stratOptions, &discordgo.ApplicationCommandOptionChoice{
			Name:  strategy.Name,
			Value: strategy.Name,
		})
	}

	draftStrategyOptions := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "draft-strategy",
			Description: "How to decide which civs to roll",
			Choices:     stratOptions,
			Required:    true,
		},
	}

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        StartDraft,
			Description: "Start a draft. You can then roll civs for players and they can submit their picks.",
			Options:     draftStrategyOptions,
		},
		{
			Name:        RollCivs,
			Description: "Roll civs for all players",
			Options:     draftStrategyOptions,
		},
		{
			Name:        Players,
			Description: "Get a list of players eligible for a draft",
		},
		{
			Name:        "ci6ndex",
			Description: "Get information about the bot",
		},
	}
	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ci6ndex":  basicCommand,
		StartDraft: startDraft,
		Players:    players,
		RollCivs:   rollCivs,
	}

	var ccommands []*discordgo.ApplicationCommand
	for _, c := range commands {
		ccmd, err := s.ApplicationCommandCreate(config.BotApplicationID, "", c)
		if err != nil {
			slog.Error("could not create slash command", "command", c.Name, "error", err)
		}
		slog.Info("registered", "command", c.Name)
		ccommands = append(ccommands, ccmd)
	}
	slog.Info("slash commands attached")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	return ccommands, nil
}

// TODO: switch to followups so we can send multiple messages
func rollCivs(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
	drafts, err := db.queries.GetActiveDrafts(context.Background())
	if err != nil {
		ReportError("error checking active drafts", err, s, i)
		return
	}

	var activeDraft domain.Ci6ndexDraft

	if len(drafts) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There is no active draft. These results will not be attached to a game",
			},
		})
		if err != nil {
			slog.Error("error responding to user", "error", err)
		}

		// dummy draft as a default
		activeDraft = domain.Ci6ndexDraft{
			ID:            -1,
			DraftStrategy: "RandomPickPool3Standard",
			Active:        true,
		}
	}

	if len(drafts) > 1 {
		ReportError("there are multiple active drafts. This should not be possible", nil, s, i)
		return
	}

	if len(drafts) == 1 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Rolled civs will be attached to the active draft",
			},
		})

		if err != nil {
			slog.Error("error responding to user", "error", err)
		}

		activeDraft = drafts[0]
	}

	picks, err := OfferPicks(db, activeDraft, 3)
	if err != nil {
		ReportError("error rolling civs", err, s, i)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Rolled civs: %v", picks),
		},
	})

	if err != nil {
		slog.Error("error responding to user", "error", err)
	}

}

func players(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)
	users, err := db.queries.GetUsers(context.Background())
	if err != nil {
		slog.Error("error getting players", "error", err)
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error getting players: %s", err.Error()),
			},
		})
		if err != nil {
			slog.Error("error responding to user", "error", err)
		}
		return
	}
	var playerNames []string
	for _, p := range users {
		playerNames = append(playerNames, p.Name)
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Players eligible for draft: %s", strings.Join(playerNames, ", ")),
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

func startDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction.ApplicationCommandData().Name)

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	strat := optionMap["draft-strategy"].StringValue()

	if strat == "" {
		slog.Error("no strategy provided for draft - this should not be possible")
		return
	}

	ds, err := db.queries.GetDraftStrategy(context.Background(), strat)

	if err != nil {
		ReportError("error fetching draft strategy", err, s, i)
		return
	}

	actives, err := db.queries.GetActiveDrafts(context.Background())

	if err != nil {
		ReportError("error fetching active draft", err, s, i)
		return
	}

	if len(actives) > 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There is already an active draft. Please end it before starting a new one.",
			},
		})
		if err != nil {
			slog.Error("error responding to user", "error", err)
		}
		return
	}

	draft, err := db.queries.CreateDraft(context.Background(), ds.Name)

	if err != nil {
		ReportError("error creating draft", err, s, i)
		return
	}

	slog.Info("draft created", "draft", draft.ID, "strategy", ds.Name, "startedBy", i.Interaction.Member.User.Username)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Draft #%v %s started by user %s. %s", draft.ID, ds.Name, i.Interaction.Member.User.Username, ds.Description),
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

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

func ready(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateGameStatus(0, "/ci6ndex")
	if err != nil {
		slog.Warn("could not update discord status on startup")
	}
	slog.Info("bot initialized and ready to receive events")
}

func ReportError(msg string, err error, s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Error(msg, "error", err)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		slog.Error("error responding to user", "error", err)
	}
}
