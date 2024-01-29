package alert

import (
	"context"
	"time"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
	"github.com/disgoorg/disgo"
	discordgo "github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
)

type discord struct {
    Client webhook.Client
    Config config.DiscordConfig
}

func NewDiscord() AlertInterface {
    return discord{}
}

func (discord discord) Start() {
	console.Normal("Starting discord client...")
    discord.Config = config.LoadDiscordAlertConfig()
	console.Normal("Discord client version: (disgo)" + disgo.Version)
	discord.Client = webhook.New(snowflake.MustParse(discord.Config.WebHookId), discord.Config.WebHookToken)
    console.Success("Discord client started! Ready to send alerts!")

	return
}

func (discord discord) Close() {
    if discord.Client != nil {
        discord.Client.Close(context.Background())
        console.Normal("Shutting down discord client...")
    }

	return
}

func (discord discord) Send(message string) {
	var _, alertError = discord.Client.CreateMessage(discordgo.NewWebhookMessageCreateBuilder().SetContent(message).Build(), rest.WithDelay(2*time.Second))
	if alertError != nil {
		console.Error("There was an error during the creation of the discord alert: " + alertError.Error())
	}

	return
}
