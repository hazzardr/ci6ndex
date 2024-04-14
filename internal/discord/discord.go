package discord

import (
	"ci6ndex/internal"
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
	mb     *MessageBuilder
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
		mb:     NewDiscTemplate(),
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
	for name, h := range getDraftHandlers(bot.db) {
		handlers[name] = h
	}
	for name, h := range getRollCivsHandlers(bot.db, bot.mb) {
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
	commands = append(commands, getRollCivCommands()...)

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

func basicCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	slog.Info("event received", "command", i.Interaction.ApplicationCommandData().Name)
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ci6ndex (Civ VI Index) is a bot for managing civ 6 draft and game information.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		slog.Error(err.Error())
	}
}

func reportError(msg string, err error, s *discordgo.Session, i *discordgo.InteractionCreate, followup bool) {
	slog.Error(msg, "interactionId", i.ID, "error", err)
	if followup {
		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong when executing the command. Interaction ID:" + i.ID,
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	} else {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong when executing the command. Interaction ID:" + i.ID,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			slog.Error("failed to respond to interaction", "i", i.Interaction.ID, "error", err)
		}
	}
}
