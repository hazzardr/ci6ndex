package internal

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log/slog"
	"time"
)

const (
	StartDraft          = "start-draft"
	FreeRollCivs        = "roll-civs"
	RollCivsForDraft    = "rollcivsdraft"
	RerollCivsForPlayer = "rollcivsplayer"
)

func freeRollCivs(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//TODO
}

func rollCivsForDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//TODO
}

func followupsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Followup messages are basically regular messages (you can create as many of them as you wish)
	// but work as they are created by webhooks and their functionality
	// is for handling additional messages after sending a response.
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			// Note: this isn't documented, but you can use that if you want to.
			// This flag just allows you to create messages visible only for the caller of the command
			// (user who triggered the command)
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Surprise!",
		},
	})
	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Followup message has been created, after 5 seconds it will be edited",
	})
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong",
		})
		return
	}
	time.Sleep(time.Second * 5)

	content := "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted."
	s.FollowupMessageEdit(i.Interaction, msg.ID, &discordgo.WebhookEdit{
		Content: &content,
	})

	time.Sleep(time.Second * 10)

	s.FollowupMessageDelete(i.Interaction, msg.ID)

	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
			"take a unicorn :unicorn: as reward!\n" +
			"Also, as bonus... look at the original interaction response :D",
	})

}

// AttachSlashCommands attaches all slash commands to the bot. Database has to be initialized first.
func AttachSlashCommands(s *discordgo.Session, config *AppConfig) {
	err := db.Health()
	if err != nil {
		panic(fmt.Errorf("can't attach commands prior to db: %w", err))
	}
	strategies, err := db.queries.GetDraftStrategies(context.Background())
	if err != nil {
		panic(fmt.Errorf("couldn't get required defaults: %w", err))
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
			Name:        FreeRollCivs,
			Description: "Roll civs for a given number of players (not associated with draft and not saved)",
			Options:     draftStrategyOptions,
		},
		{
			Name:        "ci6ndex",
			Description: "Get information about the bot",
		},
	}
	handlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ci6ndex":  basicCommand,
		StartDraft: startDraft,
	}
	for _, c := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, config.FokGuildID, c)
		if err != nil {
			slog.Error("could not create slash command", "command", c.Name, "error", err)
		}
		slog.Info("registered", "command", c.Name)
	}
	slog.Info("slash commands attached")
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

}

func startDraft(s *discordgo.Session, i *discordgo.InteractionCreate) {

}

func basicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "test test test!",
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

func optionsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("command received", "command", i.Interaction)
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs := make([]interface{}, 0, len(options))
	msgformat := "You learned how to use command options! " +
		"Take a look at the value(s) you entered:\n"

	// Get the value from the option map.
	// When the option exists, ok = true
	if option, ok := optionMap["string-option"]; ok {
		// Option values must be type asserted from interface{}.
		// Discordgo provides utility functions to make this simple.
		margs = append(margs, option.StringValue())
		msgformat += "> string-option: %s\n"
	}

	if option, ok := optionMap["string-choice"]; ok {
		margs = append(margs, option.StringValue())
		msgformat += "> string-choice: %s\n"
	}

	if opt, ok := optionMap["integer-option"]; ok {
		margs = append(margs, opt.IntValue())
		msgformat += "> integer-option: %d\n"
	}

	if opt, ok := optionMap["number-option"]; ok {
		margs = append(margs, opt.FloatValue())
		msgformat += "> number-option: %f\n"
	}

	if opt, ok := optionMap["bool-option"]; ok {
		margs = append(margs, opt.BoolValue())
		msgformat += "> bool-option: %v\n"
	}

	if opt, ok := optionMap["channel-option"]; ok {
		margs = append(margs, opt.ChannelValue(nil).ID)
		msgformat += "> channel-option: <#%s>\n"
	}

	if opt, ok := optionMap["user-option"]; ok {
		margs = append(margs, opt.UserValue(nil).ID)
		msgformat += "> user-option: <@%s>\n"
	}

	if opt, ok := optionMap["role-option"]; ok {
		margs = append(margs, opt.RoleValue(nil, "").ID)
		msgformat += "> role-option: <@&%s>\n"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

func ready(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateGameStatus(0, "!ci6ndex")
	if err != nil {
		slog.Warn("could not update discord status on startup")
	}
	slog.Info("bot initialized and ready to receive events")
}
