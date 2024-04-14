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

var (
	Ci6ndexCommand = &discordgo.ApplicationCommand{
		Name:        "ci6ndex",
		Description: "Provides information about the bot",
	}
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
		slog.Info(fmt.Sprintf("logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator))
	})

	handlers := make(map[string]CommandHandler)
	handlers[Ci6ndexCommand.Name] = basicCommand
	for name, h := range getDraftHandlers(bot.db, bot.config) {
		handlers[name] = h
	}
	bot.s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
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
	commands := make([]*discordgo.ApplicationCommand, 0)
	commands = append(commands, Ci6ndexCommand)
	commands = append(commands, getDraftCommands()...)

	for _, c := range commands {
		_, err := bot.s.ApplicationCommandCreate(bot.config.BotApplicationID, guild, c)
		if err != nil {
			slog.Error("could not create slash (/) command", "command", c.Name, "error", err)
			return nil, err
		}
		slog.Info("registered", "command", c.Name, "guildId", guild)
	}
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
	slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name)
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

func basicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ci6ndex (Civ VI Index) is a bot for managing civ 6 draft and game information.",
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

func (bot *Bot) reportError(msg string, err error, i *discordgo.InteractionCreate) {
	slog.Error(msg, "error", err)
	_, err = bot.s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Something went wrong",
	})
	if err != nil {
		slog.Error("error responding to user", "error", err)
	}
}
