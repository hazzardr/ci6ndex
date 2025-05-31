package bot

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"time"
)

func HandlePing(e *handler.CommandEvent) error {
	var gatewayPing string
	if e.Client().HasGateway() {
		gatewayPing = e.Client().Gateway.Latency().String()
	}

	eb := discord.NewEmbedBuilder().
		SetTitle("Pong!").
		AddField("Rest", "loading...", false).
		AddField("Gateway", gatewayPing, false).
		SetColor(colorSuccess)
	defer func() {
		start := time.Now().UnixNano()
		_, _ = e.Client().Rest.GetBotApplicationInfo()
		duration := time.Now().UnixNano() - start
		eb.SetField(0, "Round Trip", time.Duration(duration).String(), false)
		if _, err := e.Client().Rest.UpdateInteractionResponse(e.ApplicationID(), e.Token(), discord.MessageUpdate{Embeds: &[]discord.Embed{eb.Build()}}); err != nil {
			e.Client().Logger.Error("Failed to update ping embed: ", err)
		}
	}()
	return e.Respond(discord.InteractionResponseTypeCreateMessage, discord.NewMessageCreateBuilder().
		SetEmbeds(eb.Build()).
		Build(),
	)

}
