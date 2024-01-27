package alert

import (
	"context"
	"time"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

var discordClient webhook.Client

func InitDiscordClient() {
	console.Normal("Starting discord client...")
	config.LoadConfig("discord")
	console.Normal("Discord client version: (disgo)" + disgo.Version)
	discordClient = webhook.New(snowflake.MustParse(config.Discord.WebHookId), config.Discord.WebHookToken)
    console.Success("Discord client started! Ready to send alerts!")

	return
}

func CloseDiscordClient() {
    if discordClient != nil {
        discordClient.Close(context.Background())
        console.Normal("Shutting down discord client...")
    }

	return
}

func SendDiscordAlert(message string) {
	var _, alertError = discordClient.CreateMessage(discord.NewWebhookMessageCreateBuilder().SetContent(message).Build(), rest.WithDelay(2*time.Second))
	if alertError != nil {
		console.Error("There was an error during the creation of the discord alert: " + alertError.Error())
	}

	return
}
