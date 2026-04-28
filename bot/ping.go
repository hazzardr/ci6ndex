package bot

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func HandlePing(_ discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	var gatewayPing string
	if e.Client().HasGateway() {
		gatewayPing = e.Client().Gateway.Latency().String()
	}

	eb := discord.NewEmbed().
		WithTitle("Pong!").
		AddField("Rest", "loading...", false).
		AddField("Gateway", gatewayPing, false).
		WithColor(colorSuccess)
	defer func() {
		start := time.Now().UnixNano()
		_, _ = e.Client().Rest.GetBotApplicationInfo()
		duration := time.Now().UnixNano() - start
		eb.AddField("Round Trip", time.Duration(duration).String(), false)
		if _, err := e.Client().Rest.UpdateInteractionResponse(
			e.ApplicationID(),
			e.Token(),
			discord.MessageUpdate{Embeds: &[]discord.Embed{eb}},
		); err != nil {
			slog.Error("failed to update ping embed: ", slog.Any("err", err))
		}
	}()
	return e.Respond(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateV2().
		WithEmbeds(eb),
	)
}
